package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	sqlFalse      = "\x00"
	sqlTrue       = "\x01"
	sqlConnection = "root:banana123@tcp(127.0.0.1:3306)/loldongs?parseTime=true"
)

func getSQLConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", sqlConnection)
	return db, err
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
