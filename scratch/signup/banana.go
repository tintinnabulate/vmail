package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func insertGreeting(db *sql.DB, greeting string, awesomeness int) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO greetings(greeting, awesomeness) VALUES(?, ?)")
	checkErr(err)
	res, err := stmt.Exec(greeting, awesomeness)
	checkErr(err)
	lastId, err := res.LastInsertId()
	return lastId, err
}

func main() {
	db, err := sql.Open("mysql",
		"root:banana123@tcp(127.0.0.1:3306)/hello")
	checkErr(err)
	defer db.Close()
	lastId, err := insertGreeting(db, "yo!", 8)
	checkErr(err)
	log.Printf("ID = %d\n", lastId)
}
