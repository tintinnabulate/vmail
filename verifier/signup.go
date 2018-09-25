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
	"net/http"

    "github.com/tintinnabulate/gonfig"
	"github.com/gorilla/mux"
	"github.com/tintinnabulate/aecontext-handlers/handlers"
	"golang.org/x/net/context"
	"google.golang.org/appengine/mail"
)

// CreateHandler creates my mux.Router. Uses f to convert ContextHandlerFunc's to HandlerFunc's.
func CreateHandler(f handlers.ToHandlerHOF) *mux.Router {
	appRouter := mux.NewRouter()
	appRouter.HandleFunc("/signup/{site_code}/{email}", f(CreateSignupEndpoint)).Methods("POST")
	appRouter.HandleFunc("/verify/{site_code}/{code}", f(VerifyCodeEndpoint)).Methods("GET")
	appRouter.HandleFunc("/signup/{site_code}/{email}", f(IsSignupVerifiedEndpoint)).Methods("GET")

	return appRouter
}

// CreateSignupEndpoint handles POST /signup/{email}
func CreateSignupEndpoint(ctx context.Context, w http.ResponseWriter, req *http.Request) {
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
		codeIsOkayToUse, err := IsCodeAvailable(ctx, code)
		CheckErr(err)
		if codeIsOkayToUse {
			break
		}
	}
	if err := EmailVerificationCode(ctx, email.Address, siteCode, code); err != nil {
		email.Success = false
		email.Note = err.Error()
		err = json.NewEncoder(w).Encode(email)
		CheckErr(err)
		return
	}
	_, err := AddSignup(ctx, siteCode, email.Address, code)
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
func VerifyCodeEndpoint(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var verification Verification
	verification.Code = params["code"]
	siteCode := params["site_code"]
	verification.Success = true
	site, _ := GetSite(ctx, siteCode)
	verification.Note = site.RootURL
	err := MarkVerified(ctx, verification.Code)
	if err != nil {
		verification.Success = false
		verification.Note = err.Error()
		http.Redirect(w, req, site.RootURL, http.StatusSeeOther)
	} else {
		http.Redirect(w, req, site.VerifiedURL, http.StatusFound)
	}
}

// IsSignupVerifiedEndpoint handles GET /signup/{email}
func IsSignupVerifiedEndpoint(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	var email Email
	email.Address = params["email"]
	exists, err := IsSignupVerified(ctx, email.Address)
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

type Config struct {
	SMTPUsername string `id:"SMTPUsername" default:"sender@mydomain.com"`
	SMTPPassword string `id:"SMTPPassword" default:"mypassword"`
	SMTPServer   string `id:"SMTPServer"   default:"smtp.mydomain.com"`
	SiteDomain   string `id:"SiteDomain"   default:"mydomain.com"`
	SiteName     string `id:"SiteName"     default:"MyDomain"`
	ProjectID    string `id:"ProjectID"    default:"my-appspot-project-id"`
	ProjectURL   string `id:"ProjectURL"   default:"my-appspot-project-id.appspot.com"`
	ProjectEmail string `id:"ProjectEmail" default:"donotreply@my-appspot-project-id.appspotmail.com"`
}

// configInit : load in config file using gonfig
func configInit(configName string) {
	err := gonfig.Load(&config, gonfig.Conf{
		FileDefaultFilename: configName,
		FileDecoder:         gonfig.DecoderJSON,
		FlagDisable:         true,
	})
	checkErr(err)
}

// ComposeVerificationEmail builds the verification email, ready to be sent
func ComposeVerificationEmail(site Site, address, code string) *mail.Message {
	return &mail.Message{
		Sender:  fmt.Sprintf("[DO NOT REPLY] %s Admin <%s>", site.SiteName, config.ProjectEmail),
		To:      []string{address},
		Subject: fmt.Sprintf("[%s Registration] Please confirm your email address", site.SiteName),
		Body:    fmt.Sprintf(verificationEmailBody, site.SiteName, config.ProjectURL, site.Code, code, site.SiteName),
	}
}

// EmailVerificationCode composes and sends a verification code email
func EmailVerificationCode(ctx context.Context, address, siteCode, code string) error {
	site, _ := GetSite(ctx, siteCode)
	msg := ComposeVerificationEmail(site, address, code)
	return mail.Send(ctx, msg)
}

func init() {
	configInit("config.json")
	http.Handle("/", CreateHandler(handlers.ToHTTPHandler))
}
