package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"

	"google.golang.org/appengine"
	"google.golang.org/appengine/mail"
)

type configuration struct {
	SiteName     string
	SiteDomain   string
	SMTPServer   string
	SMTPUsername string
	SMTPPassword string
	ProjectID    string
}

var (
	config    configuration
	validPath = regexp.MustCompile("^/verify/([a-zA-Z0-9]+)$")
)

const emailBody = `
Bananas
`

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = configuration{}
	err := decoder.Decode(&config)
	checkErr(err)
	http.HandleFunc("/signup/", signupHandler)
	http.HandleFunc("/verify/", verifyHandler)
}

func randToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t, _ := template.ParseFiles("signup.html")
		t.Execute(w, nil)
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		email := string(r.FormValue("email"))
		code := randToken()
		_, err := AddSignup(r, email, code)
		if err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			fmt.Fprint(w, "Success")
		}
		ctx := appengine.NewContext(r)
		msg := &mail.Message{
			Sender:  "[DO NOT REPLY] Admin <donotreply@seraphic-lock-199316.appspotmail.com>",
			To:      []string{email},
			Subject: fmt.Sprintf("Verification Code: %s", code),
			Body:    emailBody,
		}
		if err := mail.Send(ctx, msg); err != nil {
			fmt.Fprintf(w, err.Error()) // TODO: fix. we don't want to print error to user browser
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST are supported")
	}
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	code := m[1]
	err := MarkVerified(r, code)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
	} else {
		fmt.Fprint(w, "Thank you for verifying your email address. You can now proceed with registration")
	}
}
