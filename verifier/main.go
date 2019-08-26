/*
	Implementation Note:
		None.
	Filename:
		signup.go
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/tintinnabulate/gonfig"
)

// CreateHandler creates my mux.Router. Uses f to convert ContextHandlerFunc's to HandlerFunc's.
func CreateHandler() *mux.Router {
	appRouter := mux.NewRouter()
	appRouter.HandleFunc("/signup/{site_code}/{email}", CreateSignupEndpoint).Methods("POST")
	appRouter.HandleFunc("/verify/{site_code}/{code}", VerifyCodeEndpoint).Methods("GET")
	appRouter.HandleFunc("/signup/{site_code}/{email}", IsSignupVerifiedEndpoint).Methods("GET")

	return appRouter
}

// CreateSignupEndpoint handles POST /signup/{email}
func CreateSignupEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	var email Email
	email.Address = params["email"]
	siteCode := params["site_code"]
	code := ""
	// Loop until we get a code that is available
	// TODO: handle the case where we run out of codes (or we loop forever!)
	for {
		code = RandToken()
		codeIsOkayToUse, err := IsCodeAvailable(req.Context(), code)
		CheckErr(err)
		if codeIsOkayToUse {
			break
		}
	}
	if err := EmailVerificationCode(w, req, email.Address, siteCode, code); err != nil {
		email.Success = false
		email.Note = err.Error()
		err = json.NewEncoder(w).Encode(email)
		CheckErr(err)
		return
	}
	_, err := AddSignup(req.Context(), siteCode, email.Address, code)
	if err != nil {
		email.Success = false
		email.Note = err.Error()
	} else {
		email.Success = true
	}
	err = json.NewEncoder(w).Encode(email)
	CheckErr(err)
}

// VerifyCodeEndpoint handles GET /verify/{code}
func VerifyCodeEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var verification Verification
	verification.Code = params["code"]
	siteCode := params["site_code"]
	verification.Success = true
	site, _ := GetSite(req.Context(), siteCode)
	verification.Note = site.RootURL
	err := MarkVerified(req.Context(), verification.Code)
	if err != nil {
		verification.Success = false
		verification.Note = err.Error()
		http.Redirect(w, req, site.RootURL, http.StatusNotFound)
		return
	} else {
		http.Redirect(w, req, site.VerifiedURL, http.StatusFound)
		return
	}
}

// IsSignupVerifiedEndpoint handles GET /signup/{email}
func IsSignupVerifiedEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	var email Email
	email.Address = params["email"]
	exists, err := IsSignupVerified(req.Context(), email.Address)
	if err != nil {
		email.Success = false
		email.Note = err.Error()
		err = json.NewEncoder(w).Encode(email)
		CheckErr(err)
		return
	}
	if !exists {
		email.Success = false
	} else {
		email.Success = true
	}
	err = json.NewEncoder(w).Encode(email)
	CheckErr(err)
}

var (
	appRouter mux.Router
	config    Config
)

const verificationEmailBody = `
Welcome to %s!

To get started, please click below to confirm your email address:

https://%s/verify/%s/%s

This is an automated email. Do not reply to this email address.

-- 
Best wishes,
%s Committee.
`

// Config : the configuration file format
type Config struct {
	SMTPUsername      string `id:"SMTPUsername"      default:"sender@mydomain.com"`
	SMTPPassword      string `id:"SMTPPassword"      default:"mypassword"`
	SMTPServer        string `id:"SMTPServer"        default:"smtp.mydomain.com"`
	SiteDomain        string `id:"SiteDomain"        default:"mydomain.com"`
	SiteName          string `id:"SiteName"          default:"MyDomain"`
	ProjectID         string `id:"ProjectID"         default:"my-appspot-project-id"`
	ProjectURL        string `id:"ProjectURL"        default:"my-appspot-project-id.appspot.com"`
	ProjectEmail      string `id:"ProjectEmail"      default:"donotreply@my-appspot-project-id.appspotmail.com"`
	SendGridKey       string `id:"SendGridKey"       default:"SendGridKey"`
	GoogleCredentials string `id:"GoogleCredentials" default:"GoogleCredentialsJSONFilename"`
}

// configInit : load in config file using gonfig
func configInit(configName string) {
	err := gonfig.Load(&config, gonfig.Conf{
		FileDefaultFilename: configName,
		FileDecoder:         gonfig.DecoderJSON,
		FlagDisable:         true,
	})
	if err != nil {
		log.Fatalf("could not load configuration file: %v", err)
	}
}

// environmentInit : set up environment variables
func environmentInit() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.GoogleCredentials)
}

// ComposeVerificationEmail builds the verification email, ready to be sent
func ComposeVerificationEmail(site Site, address, code string) []byte {
	m := mail.NewV3Mail()

	sender_address := config.ProjectEmail
	sender_name := fmt.Sprintf("[DO NOT REPLY] %v Admin", site.SiteName)
	e := mail.NewEmail(sender_name, sender_address)
	m.SetFrom(e)

	subject := fmt.Sprintf("[%s Registration] Please confirm your email address", site.SiteName)
	m.Subject = subject

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(address, address),
	}
	p.AddTos(tos...)
	m.AddPersonalizations(p)

	plainTextContent := fmt.Sprintf(verificationEmailBody, site.SiteName, config.ProjectURL, site.Code, code, site.SiteName)
	c := mail.NewContent("text/plain", plainTextContent)
	m.AddContent(c)
	c = mail.NewContent("text/html", plainTextContent)
	m.AddContent(c)

	mailSettings := mail.NewMailSettings()
	bypassListManagement := mail.NewSetting(true)
	mailSettings.SetBypassListManagement(bypassListManagement)
	spamCheckSetting := mail.NewSpamCheckSetting()
	spamCheckSetting.SetEnable(true)
	spamCheckSetting.SetSpamThreshold(1)
	spamCheckSetting.SetPostToURL("https://spamcatcher.sendgrid.com")
	mailSettings.SetSpamCheckSettings(spamCheckSetting)
	m.SetMailSettings(mailSettings)

	trackingSettings := mail.NewTrackingSettings()
	clickTrackingSettings := mail.NewClickTrackingSetting()
	clickTrackingSettings.SetEnable(true)
	clickTrackingSettings.SetEnableText(true)
	trackingSettings.SetClickTracking(clickTrackingSettings)
	openTrackingSetting := mail.NewOpenTrackingSetting()
	openTrackingSetting.SetEnable(true)
	trackingSettings.SetOpenTracking(openTrackingSetting)
	subscriptionTrackingSetting := mail.NewSubscriptionTrackingSetting()
	subscriptionTrackingSetting.SetEnable(true)
	trackingSettings.SetSubscriptionTracking(subscriptionTrackingSetting)
	m.SetTrackingSettings(trackingSettings)

	return mail.GetRequestBody(m)
}

// EmailVerificationCode composes and sends a verification code email
func EmailVerificationCode(w http.ResponseWriter, req *http.Request, address, siteCode, code string) error {
	site, _ := GetSite(req.Context(), siteCode)
	request := sendgrid.GetRequest(config.SendGridKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	// attach all the QR codes into a message bound for that email address
	request.Body = ComposeVerificationEmail(site, address, code)
	// send the email
	_, err := sendgrid.API(request)
	if err != nil {
		return err
	}
	return nil

}

func init() {
	configInit("config.json")
	environmentInit()
	http.Handle("/", CreateHandler())
}

// main : main entry point to application
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Start here: http://localhost:%s/signup", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
