package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// getDB opens a database connection with file dbName.db
func getDB(dbFilename string) *sql.DB {
	db, err := sql.Open("sqlite3", dbFilename)
	checkErr(err)
	return db
}

func (db *sql.DB) insertUser(username string) {
	// insert
	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	checkErr(err)

	res, err := stmt.Exec("simon", "devteam", "2018-04-04")
	checkErr(err)
}

func (db *sql.DB) updateUsername(id int, username string) int {
	// update
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err = stmt.Exec(username, id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)
	return affected
}

func (db *sql.DB) getAllUsers() *sql.Rows {
	// query
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)
	return rows
}

func (rs *sql.Rows) printUsers() {

	var uid int
	var username string
	var department string
	var created time.Time

	for rs.Next() {
		err = rs.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println("uid:", uid)
		fmt.Println("username:", username)
		fmt.Println("department:", department)
		fmt.Println("created:", created.String())
	}
	rs.Close()
}

func (db *sql.DB) deleteAllUser() int {
	// delete
	res, err = db.Exec("DELETE FROM userinfo")
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)
	return affect

	// db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
