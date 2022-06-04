package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/philcantcode/localmapper/utils"
)

var Con *sql.DB

func InitSqlite() {
	utils.Log("Attempting to open SQL database: "+utils.Configs["DATABASE_PATH"], true)
	var err error
	Con, err = sql.Open("sqlite3", utils.Configs["DATABASE_PATH"])

	utils.ErrorLog("Couldn't open SQLite3 database", err, true)

	stmt, err := Con.Prepare(
		"CREATE TABLE IF NOT EXISTS HostTracking" +
			"(id INTEGER PRIMARY KEY UNIQUE, " +
			"mac TEXT, " +
			"ipv4 TEXT, " +
			"ipv6 TEXT, " +
			"metadata TEXT, " +
			"firstSeen DATETIME, " +
			"lastSeen DATETIME)")
	utils.ErrorLog("Couldn't create SQL database HostTracking", err, true)
	stmt.Exec()

	stmt, err = Con.Prepare(
		"CREATE TABLE IF NOT EXISTS Capabilities" +
			"(id INTEGER PRIMARY KEY UNIQUE, " +
			"name TEXT NOT NULL, " +
			"type TEXT NOT NULL, " +
			"command TEXT NOT NULL, " +
			"description TEXT DEFAULT '', " +
			"displayFields TEXT DEFAULT '')")
	utils.ErrorLog("Couldn't create SQL database Capabilities", err, true)
	stmt.Exec()

	stmt, err = Con.Prepare(
		"CREATE TABLE IF NOT EXISTS Vlans" +
			"(id INTEGER PRIMARY KEY UNIQUE, " +
			"name TEXT NOT NULL Unique, " +
			"description TEXT NOT NULL, " +
			"highIP TEXT NOT NULL, " +
			"lowIP TEXT NOT NULL)")
	utils.ErrorLog("Couldn't create SQL database Vlans", err, true)
	stmt.Exec()

	stmt, err = Con.Prepare(
		"CREATE TABLE IF NOT EXISTS Propositions" +
			"(id INTEGER PRIMARY KEY UNIQUE, " +
			"type TEXT, " +
			"date DATETIME DEFAULT CURRENT_TIMESTAMP, " +
			"description TEXT, " +
			"proposition TEXT, " +
			"correction TEXT, " +
			"status INTEGER default 0, " +
			"user INTEGER)")
	utils.ErrorLog("Couldn't create SQL database Propositions", err, true)
	stmt.Exec()

	stmt, err = Con.Prepare(
		"CREATE TABLE IF NOT EXISTS CMDB" +
			"(id INTEGER PRIMARY KEY UNIQUE, " +
			"osiLayer INTEGER NOT NULL, " +
			"dateSeen TEXT NOT NULL, " +
			"description TEXT NOT NULL, " +
			"statusTags TEXT NOT NULL, " +
			"userTags TEXT NOT NULL, " +
			"infoTags TEXT NOT NULL)")
	utils.ErrorLog("Couldn't create SQL database CMDB", err, true)
	stmt.Exec()

	stmt, err = Con.Prepare(
		"CREATE TABLE IF NOT EXISTS JobSpecs" +
			"(id INTEGER PRIMARY KEY UNIQUE, " +
			"job TEXT NOT NULL)")
	utils.ErrorLog("Couldn't create SQL database JobSpecs", err, true)
	stmt.Exec()

	stmt.Close()
}
