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
	"google.golang.org/appengine"
	"google.golang.org/appengine/mail"
)

// CreateHandler creates my mux.Router. Uses f to convert ContextHandlerFunc's to HandlerFunc's.
func CreateHandler(f ContextHandlerToHandlerHOF) *mux.Router {
	appRouter := mux.NewRouter()
	appRouter.HandleFunc("/signup/{email}", f(CreateSignupEndpoint)).Methods("POST")
	appRouter.HandleFunc("/verify/{code}", f(VerifyCodeEndpoint)).Methods("GET")
	appRouter.HandleFunc("/signup/{email}", f(IsSignupVerifiedEndpoint)).Methods("GET")

	return appRouter
}

// CreateSignupEndpoint handles POST /signup/{email}
func CreateSignupEndpoint(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	var email Email
	email.Address = params["email"]
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
	if err := EmailVerificationCode(ctx, email.Address, code); err != nil {
		email.Success = false
		email.Note = err.Error()
		json.NewEncoder(w).Encode(email)
		return
	}
	_, err := AddSignup(ctx, email.Address, code)
	if err != nil {
		email.Success = false
		email.Note = err.Error()
	} else {
		email.Success = true
	}
	json.NewEncoder(w).Encode(email)
}

// VerifyCodeEndpoint handles GET /verify/{code}
func VerifyCodeEndpoint(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	var verification Verification
	verification.Code = params["code"]
	err := MarkVerified(ctx, verification.Code)
	verification.Success = true
	verification.Note = ""
	if err != nil {
		verification.Success = false
		verification.Note = err.Error()
	}
	json.NewEncoder(w).Encode(verification)
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
		json.NewEncoder(w).Encode(email)
		return
	}
	if !exists {
		email.Success = false
	} else {
		email.Success = true
	}
	json.NewEncoder(w).Encode(email)
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

// Email holds our JSON response for GET and POST /signup/{email}
type Email struct {
	Address string `json:"address"`
	Success bool   `json:"success"`
	Note    string `json:"note"`
}

// Verification holds our JSON response for GET /verify/{code}
type Verification struct {
	Code    string `json:"code"`
	Success bool   `json:"success"`
	Note    string `json:"note"`
}

var (
	config    configuration
	appRouter mux.Router
)

const verificationEmailBody = `
Welcome to %s!

To get started, please click below to confirm your email address:

https://%s/verify/%s

Best wishes,
%s Committee.
`

// LoadConfig loads the app configuration JSON into the `config` variable
func LoadConfig() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = configuration{}
	err := decoder.Decode(&config)
	CheckErr(err)
}

// ComposeVerificationEmail builds the verification email, ready to be sent
func ComposeVerificationEmail(address, code string) *mail.Message {
	return &mail.Message{
		Sender:  fmt.Sprintf("[DO NOT REPLY] Admin <%s>", config.ProjectEmail),
		To:      []string{address},
		Subject: fmt.Sprintf("[%s] Please confirm your account", config.SiteName),
		Body:    fmt.Sprintf(verificationEmailBody, config.SiteName, config.ProjectURL, code, config.SiteName),
	}
}

// EmailVerificationCode composes and sends a verification code email
func EmailVerificationCode(ctx context.Context, address, code string) error {
	msg := ComposeVerificationEmail(address, code)
	return mail.Send(ctx, msg)
}

func init() {
	LoadConfig()
	http.Handle("/", CreateHandler(ContextHandlerToHttpHandler))
}

// HandlerFunc is our Standard http handler
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// ContextHandlerFunc is our context.Context http handler
type ContextHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)

// ContextHandlerToHandlerHOF is our Higher order function
// for changing a HandlerFunc to a ContextHandlerFunc, usually creating the context.Context along the way.
type ContextHandlerToHandlerHOF func(f ContextHandlerFunc) HandlerFunc

// ContextHandlerToHttpHandler Creates a new Context and uses it when calling f
func ContextHandlerToHttpHandler(f ContextHandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		f(ctx, w, r)
	}
}
