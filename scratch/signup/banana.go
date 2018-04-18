package main

// Got up to here: http://go-database-sql.org/modifying.html

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

type Greeting struct {
	greeting    string
	awesomeness int
}

func (g Greeting) insert(db *sql.DB) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO greetings(greeting, awesomeness) VALUES(?, ?)")
	checkErr(err)
	res, err := stmt.Exec(g.greeting, g.awesomeness)
	checkErr(err)
	lastId, err := res.LastInsertId()
	return lastId, err
}

func fetchGreeting(db *sql.DB, id int64) (Greeting, error) {
	var g Greeting
	err := db.QueryRow("SELECT greeting, awesomeness FROM greetings WHERE id = ?", id).Scan(&g.greeting, &g.awesomeness)
	return g, err
}

func deleteGreeting(db *sql.DB, id int64) error {
	stmt, err := db.Prepare("DELETE FROM greetings WHERE id = ?")
	checkErr(err)
	_, err = stmt.Exec(id)
	return err
}

func main() {
	db, err := sql.Open("mysql",
		"root:banana123@tcp(127.0.0.1:3306)/hello")
	checkErr(err)
	defer db.Close()
	g := &Greeting{"hello", 7}
	lastId, err := g.insert(db)
	checkErr(err)
	log.Printf("ID = %d\n", lastId)
	thatG, err := fetchGreeting(db, lastId)
	checkErr(err)
	log.Printf("%+v\n", thatG)
	err = deleteGreeting(db, lastId)
	checkErr(err)
}
