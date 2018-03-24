package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type Configuration struct {
	SMTPServer   string
	SMTPUsername string
	SMTPPassword string
}

var config Configuration

func init() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}
}

type Signup struct {
	Email          []byte
	ValidationCode string
}

func (s *Signup) save() error {
	filename := string(s.Email) + ".txt"
	return ioutil.WriteFile(filename, []byte(s.ValidationCode), 0600)
}

func randToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func emailCode(recipient, code string) error {
	// Set up authentication information.
	// Port and encryption:
	// - 587 with STARTTLS (recommended)
	// - 465 with TLS
	// - 25 with STARTTLS or none
	// Authentication: your email address and password

	fmt.Println(config.SMTPPassword)

	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPServer)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{recipient}
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: discount Gophers!\r\n"+
		"\r\n"+
		"Code: %s\r\n", recipient, config.SMTPUsername, code))
	return smtp.SendMail(fmt.Sprintf("%s:25", config.SMTPServer), auth, config.SMTPUsername, to, msg)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("signup.html")
	s := &Signup{}
	t.Execute(w, s)
}

func signupSubmitHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	s := &Signup{Email: []byte(email), ValidationCode: randToken()}
	err := s.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err2 := emailCode(string(s.Email), s.ValidationCode)
	if err2 != nil {
		log.Fatal(err2)
	}
	http.Redirect(w, r, "/signup/", http.StatusFound)
}

func main() {
	http.HandleFunc("/signup/", signupHandler)
	http.HandleFunc("/signup_submit/", signupSubmitHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
