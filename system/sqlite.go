package system

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Con *sql.DB

func InitSqlite() {
	fmt.Println("Attempting to open SQL database: res/database.db")

	var err error
	Con, err = sql.Open("sqlite3", "res/database.db")

	if err != nil {
		fmt.Println("Couldn't open SQLite3 database")
		fmt.Println(err)
		os.Exit(0)
	}

	stmt, err := Con.Prepare(
		"CREATE TABLE IF NOT EXISTS Settings" +
			"(id INTEGER PRIMARY KEY UNIQUE, " +
			"key TEXT UNIQUE, " +
			"value TEXT)")

	if err != nil {
		fmt.Println("Couldn't create SQL database Settings")
		fmt.Println(err)
		os.Exit(0)
	}

	_, err = stmt.Exec()

	if err != nil {
		fmt.Println("Couldn't execute SQL database Settings")
		fmt.Println(err)
		os.Exit(0)
	}

	stmt, err = Con.Prepare(
		"CREATE TABLE IF NOT EXISTS Wordlists" +
			"(id INTEGER PRIMARY KEY UNIQUE, " +
			"label TEXT UNIQUE, " +
			"description TEXT, " +
			"path TEXT, " +
			"type TEXT, " +
			"size INT)")

	if err != nil {
		fmt.Println("Couldn't create SQL database Wordlists")
		fmt.Println(err)
		os.Exit(0)
	}

	_, err = stmt.Exec()

	if err != nil {
		fmt.Println("Couldn't execute SQL database Wordlists")
		fmt.Println(err)
		os.Exit(0)
	}

	stmt.Close()
}
