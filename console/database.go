package console

import (
	"database/sql"

	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
)

func CheckCommandExists(command string, params string) int {
	utils.Log("Querying Hosts from CommandCapability DB", false)
	stmt, err := database.Connection.Prepare("SELECT `id` FROM `CommandCapability` WHERE" +
		"`command` = ? AND `params` = ?;")
	utils.ErrorHandle("Couldn't select id from CommandCapability CheckCommandExists", err, true)

	row := stmt.QueryRow(command, params)
	utils.ErrorHandle("Couldn't recieve rows from CommandCapability CheckCommandExists", err, true)

	var id int
	err = row.Scan(&id)

	if err == sql.ErrNoRows {
		return -1
	}

	return id
}

func InsertCommand(command string, params string, name string, cmdType string, desc string) {
	utils.Log("Inserting Hosts from CommandCapability DB", false)
	stmt, err := database.Connection.Prepare("INSERT INTO `CommandCapability`" +
		"(`command`, `params`, `type`, `name`, `description`) VALUES (?, ?, ?, ?, ?);")
	utils.ErrorHandle("Couldn't prepare InsertCommand CommandCapability", err, true)

	_, err = stmt.Exec(command, params, cmdType, name, desc)
	utils.ErrorHandle("Error executing CommandCapability insertHosts", err, true)
	stmt.Close()
}

func UpdateCommandDisplayField(displayFields string, id int) {
	utils.Log("Updating Hosts from CommandCapability DB", false)
	stmt, err := database.Connection.Prepare("UPDATE `CommandCapability` SET `dispalyFields` = ? WHERE `id` = ?;")
	utils.ErrorHandle("Couldn't update CommandCapability", err, true)

	_, err = stmt.Exec(displayFields, id)
	utils.ErrorHandle("Results error from CommandCapability", err, true)
	stmt.Close()
}
