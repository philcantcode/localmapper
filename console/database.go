package console

import (
	"database/sql"
	"encoding/json"
	"fmt"

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

func InsertCapability(command string, params string, name string, cmdType string, desc string) {
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

type Capabilities struct {
	id            int
	command       string
	params        []string
	cmdType       string
	name          string
	desc          string
	displayFields string
	interpreter   string
}

func GetAllCapabilities() {
	utils.Log("Querying capabilities from CommandCapability DB", false)
	stmt, err := database.Connection.Prepare("SELECT `id`, `command`, `params`, `type`, `name`, `description`, `displayFields`, `interpreter` FROM `CommandCapability`")
	utils.ErrorHandle("Couldn't select all from CommandCapability GetAllCapabilities", err, true)

	rows, err := stmt.Query()
	defer rows.Close()
	utils.ErrorHandle("Couldn't recieve rows from CommandCapability GetAllCapabilities", err, true)

	capabilities := []Capabilities{}

	for rows.Next() {
		capability := Capabilities{}
		params := ""
		var paramArr []string

		rows.Scan(capability.id, capability.command, params, capability.cmdType, capability.name, capability.desc, capability.displayFields, capability.interpreter)
		json.Unmarshal([]byte(params), &paramArr)
		capability.params = paramArr

		capabilities = append(capabilities, capability)

		fmt.Printf("+ %v\n", capability)
	}
}
