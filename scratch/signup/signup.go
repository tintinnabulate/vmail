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
	"regexp"
)

type configuration struct {
	SiteName     string
	SiteDomain   string
	SMTPServer   string
	SMTPUsername string
	SMTPPassword string
}

var config configuration

var validPath = regexp.MustCompile("^/verify/([a-zA-Z0-9]+)$")

func init() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}
}

type signup struct {
	Email            []byte
	VerificationCode string
}

func (s *signup) save() error {
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
	s := &signup{}
	t.Execute(w, s)
}

func signupSubmitHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	s := &signup{Email: []byte(email), VerificationCode: randToken()}
	err := s.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	emailCode(string(s.Email), s.VerificationCode)
	http.Redirect(w, r, "/signup/", http.StatusFound)
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	fmt.Println(m)
	code := m[1]
	fmt.Fprint(w, code)
}

func main() {
	http.HandleFunc("/signup/", signupHandler)
	http.HandleFunc("/signup_submit/", signupSubmitHandler)
	http.HandleFunc("/verify/", verifyHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
