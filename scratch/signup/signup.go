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
	"os"

	"github.com/gorilla/mux"

	"golang.org/x/net/context"
	"google.golang.org/appengine/mail"
)

// CreateHandler creates my mux.Router. Uses f to convert ContextHandlerFunc's to HandlerFunc's.
func CreateHandler(f ContextHandlerToHandlerHOF) *mux.Router {
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
	}
	http.Redirect(w, req, site.VerifiedURL, http.StatusFound)
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

// configuration holds our app configuration settings
type configuration struct {
	SiteName     string
	SiteDomain   string
	SMTPServer   string
	SMTPUsername string
	SMTPPassword string
	ProjectID    string
	ProjectEmail string
	ProjectURL   string
}

var (
	config    configuration
	appRouter mux.Router
)

const verificationEmailBody = `
Welcome to %s!

To get started, please click below to confirm your email address:

https://%s/verify/%s/%s

-- 
Best wishes,
%s Committee.
`

// LoadConfig loads the app configuration JSON into the `config` variable
func LoadConfig() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config = configuration{}
	err := decoder.Decode(&config)
	CheckErr(err)
	err = file.Close()
	CheckErr(err)
}

// ComposeVerificationEmail builds the verification email, ready to be sent
func ComposeVerificationEmail(address, siteCode, code string) *mail.Message {
	return &mail.Message{
		Sender:  fmt.Sprintf("[DO NOT REPLY] Admin <%s>", config.ProjectEmail),
		To:      []string{address},
		Subject: fmt.Sprintf("[%s] Please confirm your account", config.SiteName),
		Body:    fmt.Sprintf(verificationEmailBody, config.SiteName, config.ProjectURL, siteCode, code, config.SiteName),
	}
}

// EmailVerificationCode composes and sends a verification code email
func EmailVerificationCode(ctx context.Context, address, siteCode, code string) error {
	msg := ComposeVerificationEmail(address, siteCode, code)
	return mail.Send(ctx, msg)
}

func init() {
	LoadConfig()
	http.Handle("/", CreateHandler(ContextHandlerToHTTPHandler))
}
