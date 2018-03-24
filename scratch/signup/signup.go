package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Signup struct {
	Email []byte
}

func (s *Signup) save() error {
	filename := string(s.Email) + ".txt"
	fmt.Println("hoi!")
	return ioutil.WriteFile(filename, s.Email, 0600)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("signup.html")
	s := &Signup{}
	t.Execute(w, s)
}

func signupSubmitHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	s := &Signup{Email: []byte(email)}
	err := s.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/signup/", http.StatusFound)
}

func main() {
	http.HandleFunc("/signup/", signupHandler)
	http.HandleFunc("/signup_submit/", signupSubmitHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
