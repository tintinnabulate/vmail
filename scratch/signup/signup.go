package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
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
