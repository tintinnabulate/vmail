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

func Monkeys(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	fmt.Fprintf(w, "banana: %s", params["poop"])
}

func VerifyCodeEndpoint(ctx context.Context, w http.ResponseWriter, req *http.Request) {
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

func CreateSignupEndpoint(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	var email Email
	fmt.Println("got here!")
	fmt.Printf("body: %v", req.Body)
	err := json.NewDecoder(req.Body).Decode(&email)
	fmt.Println("got here 2!")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	email.Address = params["email"]
	fmt.Println("got here 3!")
	code := randToken()
	fmt.Println("got here 4!")
	if err := EmailVerificationCode(ctx, email.Address, code); err != nil {
		email.Success = false
		email.Note = err.Error()
		json.NewEncoder(w).Encode(email)
		return
	}
	fmt.Println("got here 5!")
	_, err = AddSignup(req, email.Address, code)
	fmt.Println("got here 6!")
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

func LoadConfig() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = configuration{}
	err := decoder.Decode(&config)
	checkErr(err)
}

func randToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

/*
Standard http handler
*/
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

/*
Our context.Context http handler
*/
type ContextHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)

/*
  Higher order function for changing a HandlerFunc to a ContextHandlerFunc,
  usually creating the context.Context along the way.
*/
type ContextHandlerToHandlerHOF func(f ContextHandlerFunc) HandlerFunc

/*
Creates a new Context and uses it when calling f
*/
func ContextHanderToHttpHandler(f ContextHandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		f(ctx, w, r)
	}
}

/*
Creates my mux.Router. Uses f to convert ContextHandlerFunc's to HandlerFuncs.
*/
func CreateHandler(f ContextHandlerToHandlerHOF) *mux.Router {
	appRouter := mux.NewRouter()
	appRouter.HandleFunc("/verify/{code}", f(VerifyCodeEndpoint)).Methods("GET")
	appRouter.HandleFunc("/signup/{email}", f(CreateSignupEndpoint)).Methods("POST")
	appRouter.HandleFunc("/monkeys/{poop}", f(Monkeys)).Methods("GET")

	return appRouter
}

func init() {
	LoadConfig()
	http.Handle("/", CreateHandler(ContextHanderToHttpHandler))
}
