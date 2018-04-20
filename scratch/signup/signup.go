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
	"time"
)

type configuration struct {
	SiteName     string
	SiteDomain   string
	SMTPServer   string
	SMTPUsername string
	SMTPPassword string
}

var (
	config    configuration
	validPath = regexp.MustCompile("^/verify/([a-zA-Z0-9]+)$")
)

func init() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = configuration{}
	err := decoder.Decode(&config)
	checkErr(err)
	http.HandleFunc("/signup/", signupHandler)
	http.HandleFunc("/signup_submit/", signupSubmitHandler)
	http.HandleFunc("/verify/", verifyHandler)
}

type signup struct {
	ID                int64
	CreationTimestamp time.Time
	Email             []byte
	VerificationCode  string
	IsVerified        []uint8
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
	s := &signup{
		Email:            []byte(email),
		VerificationCode: randToken(),
		IsVerified:       []byte(sqlFalse)}
	db, err := getSQLConnection()
	checkErr(err)
	defer db.Close()
	_, err = s.insert(db)
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
	db, err := getSQLConnection()
	checkErr(err)
	defer db.Close()
	isValid, reason, maybeSignup := isValidVerificationCode(db, code)
	if isValid {
		err := maybeSignup.verify(db)
		checkErr(err)
		fmt.Fprint(w, reason)
	} else {
		fmt.Fprint(w, reason)
	}
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}
