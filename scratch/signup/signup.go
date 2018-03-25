package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
)

type Configuration struct {
	SiteName     string
	SiteDomain   string
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
	Email            []byte
	VerificationCode string
}

func (s *Signup) save() error {
	filename := string(s.Email) + ".txt"
	return ioutil.WriteFile(filename, []byte(s.VerificationCode), 0600)
}

func randToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func emailCode(recipient, code string) {
	server := config.SMTPServer
	from := mail.Address{"", config.SMTPUsername}
	to := mail.Address{"", recipient}
	subj := fmt.Sprintf("%s Verification Code", config.SiteName)
	body := fmt.Sprintf("Code: %s", code)
	sendEmail(server, from, to, subj, body)
}

func makeEmailHeaders(from, to mail.Address, subject string) map[string]string {
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	return headers
}

func makeEmailMessage(from, to mail.Address, subject, body string) string {
	message := ""
	for k, v := range makeEmailHeaders(from, to, subject) {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	return message
}

func sendEmail(server string, from, to mail.Address, subject, body string) {
	message := makeEmailMessage(from, to, subject, body)

	// Connect to the SMTP Server
	serverName := fmt.Sprintf("%s:465", server)
	host, _, _ := net.SplitHostPort(serverName)
	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPServer)

	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// For smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", serverName, tlsConfig)
	if err != nil {
		log.Panic(err)
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}
	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}
	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}
	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}
	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}
	err = w.Close()
	if err != nil {
		log.Panic(err)
	}
	c.Quit()
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("signup.html")
	s := &Signup{}
	t.Execute(w, s)
}

func signupSubmitHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	s := &Signup{Email: []byte(email), VerificationCode: randToken()}
	err := s.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	emailCode(string(s.Email), s.VerificationCode)
	http.Redirect(w, r, "/signup/", http.StatusFound)
}

func main() {
	http.HandleFunc("/signup/", signupHandler)
	http.HandleFunc("/signup_submit/", signupSubmitHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
