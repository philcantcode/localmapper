package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/philcantcode/localmapper/utils"
)

var connection *sql.DB

func InitSqlite() {
	utils.Log("Attempting to open SQL database: "+utils.Configs["DATABASE_PATH"], true)
	var err error
	connection, err = sql.Open("sqlite3", utils.Configs["DATABASE_PATH"])

	utils.ErrorLog("Couldn't open SQLite3 database", err, true)

	stmt, err := connection.Prepare(
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

	stmt, err = connection.Prepare(
		"CREATE TABLE IF NOT EXISTS Capabilities" +
			"(id INTEGER PRIMARY KEY UNIQUE, " +
			"name TEXT NOT NULL, " +
			"type TEXT NOT NULL, " +
			"params TEXT NOT NULL, " +
			"description TEXT DEFAULT '', " +
			"displayFields TEXT DEFAULT '')")
	utils.ErrorLog("Couldn't create SQL database NmapScripts", err, true)
	stmt.Exec()

	stmt.Close()
}
