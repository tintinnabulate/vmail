package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/mail"
)

func VerifyCodeEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	var verification Verification
	_ = json.NewDecoder(req.Body).Decode(&verification)
	verification.Code = params["code"]
	err := MarkVerified(req, verification.Code)
	verification.Success = true
	verification.Note = ""
	if err != nil {
		verification.Success = false
		verification.Note = err.Error()
	}
	json.NewEncoder(w).Encode(verification)
}

func ComposeVerificationEmail(address, code string) *mail.Message {
	return &mail.Message{
		Sender:  "[DONUT REPLY] Admin <donotreply@seraphic-lock-199316.appspotmail.com>",
		To:      []string{address},
		Subject: "Your verification code",
		Body:    fmt.Sprintf(emailBody, code),
	}
}

func EmailVerificationCode(ctx context.Context, address, code string) error {
	msg := ComposeVerificationEmail(address, code)
	return mail.Send(ctx, msg)
}

func CreateSignupEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	var email Email
	_ = json.NewDecoder(req.Body).Decode(&email)
	email.Address = params["email"]
	code := randToken()
	ctx := appengine.NewContext(req)
	if err := EmailVerificationCode(ctx, email.Address, code); err != nil {
		email.Success = false
		email.Note = err.Error()
		json.NewEncoder(w).Encode(email)
		return
	}
	_, err := AddSignup(req, email.Address, code)
	if err != nil {
		email.Success = false
		email.Note = err.Error()
	} else {
		email.Success = true
	}
	json.NewEncoder(w).Encode(email)
}

type configuration struct {
	SiteName     string
	SiteDomain   string
	SMTPServer   string
	SMTPUsername string
	SMTPPassword string
	ProjectID    string
}

type Email struct {
	Address string `json:"address"`
	Success bool   `json:"success"`
	Note    string `json:"note"`
}

type Verification struct {
	Code    string `json:"code"`
	Success bool   `json:success"`
	Note    string `json:note"`
}

var (
	config    configuration
	appRouter mux.Router
)

const emailBody = `
Code: %s

Yours randomly,
Bert.
`

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Initialise() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = configuration{}
	err := decoder.Decode(&config)
	checkErr(err)
	appRouter := mux.NewRouter()
	appRouter.HandleFunc("/verify/{code}", VerifyCodeEndpoint).Methods("GET")
	appRouter.HandleFunc("/signup/{email}", CreateSignupEndpoint).Methods("POST")
	http.Handle("/", appRouter)
}

func init() {
	Initialise()
}

func randToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}
