package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
	db        *sql.DB
	validPath = regexp.MustCompile("^/verify/([a-zA-Z0-9]+)$")
	LOCATION  = time.UTC
	sqlFalse  = "\x00"
	sqlTrue   = "\x01"
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

// Datetime is a utility function for making dates with the same
// location, 0 hours, 0 mins, 0 secs, 0 nanosecs.
func Datetime(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, LOCATION)
}

// Now returns the time right now when it is called
func Now() time.Time {
	return time.Now().In(LOCATION)
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
	db, err := sql.Open("mysql",
		"root:banana123@tcp(127.0.0.1:3306)/loldongs?parseTime=true")
	checkErr(err)
	defer db.Close()
	lastId, err := s.insert(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(lastId)
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
	db, err := sql.Open("mysql",
		"root:banana123@tcp(127.0.0.1:3306)/loldongs?parseTime=true")
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

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (s *signup) insert(db *sql.DB) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO " +
		"signups(email, verification_code, is_verified) " +
		"VALUES(?, ?, ?);")
	checkErr(err)
	res, err := stmt.Exec(
		s.Email,
		s.VerificationCode,
		s.IsVerified)
	checkErr(err)
	lastId, err := res.LastInsertId()
	return lastId, err
}

func isValidVerificationCode(db *sql.DB, code string) (bool, string, signup) {
	// get signup using code
	stmt, err := db.Prepare(
		`SELECT 
		id, creation_timestamp, email, 
		verification_code, is_verified 
		FROM signups WHERE verification_code = ?;`)
	checkErr(err)
	defer stmt.Close()
	var maybeSignup signup
	// populate signup
	err = stmt.QueryRow(code).Scan(
		&maybeSignup.ID,
		&maybeSignup.CreationTimestamp,
		&maybeSignup.Email,
		&maybeSignup.VerificationCode,
		&maybeSignup.IsVerified)
	switch {
	// no match
	case err == sql.ErrNoRows:
		return false,
			"No such verification code",
			maybeSignup
	case err != nil:
		log.Fatal(err)
		return false,
			"",
			maybeSignup
	// match
	default:
		// already verified, needs to continue
		if string(maybeSignup.IsVerified) == sqlTrue {
			return false,
				"You are already verified. " +
					"Please continue with the Registration process: /signup",
				maybeSignup
		} else {
			return true,
				"You are now verified. " +
					"Please continue with the Registration process: /signup",
				maybeSignup
		}
	}
}

func (s *signup) verify(db *sql.DB) error {
	stmt2, err := db.Prepare(
		`UPDATE signups
		SET is_verified = true
	    WHERE id = ? 
		AND is_verified = false;`)
	defer stmt2.Close()
	_, err = stmt2.Exec(s.ID)
	return err
}

// TODO fix query - not enough args
func fetchSignup(db *sql.DB, id int64) (signup, error) {
	var s signup
	err := db.QueryRow("SELECT creation_timestamp, email, verification_code FROM signups WHERE id = ?", id).Scan(
		&s.CreationTimestamp,
		&s.Email,
		&s.VerificationCode)
	return s, err
}

func deleteSignup(db *sql.DB, id int64) error {
	stmt, err := db.Prepare("DELETE FROM signups WHERE id = ?")
	checkErr(err)
	_, err = stmt.Exec(id)
	return err
}
