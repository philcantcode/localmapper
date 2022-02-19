package execute

import (
	"database/sql"
	"encoding/json"

	"github.com/philcantcode/localmapper/apps/database"
	"github.com/philcantcode/localmapper/utils"
)

func CheckCommandExists(params string) int {
	utils.Log("Querying Hosts from CommandCapability DB", false)
	stmt, err := database.Connection.Prepare("SELECT `id` FROM `CommandCapability` WHERE" +
		"`cmdParams` = ?;")
	utils.ErrorLog("Couldn't select id from CommandCapability CheckCommandExists", err, true)

	row := stmt.QueryRow(params)
	utils.ErrorLog("Couldn't recieve rows from CommandCapability CheckCommandExists", err, true)

	var id int
	err = row.Scan(&id)

	if err == sql.ErrNoRows {
		return -1
	}

	return id
}

func InsertCapability(params string, name string, cmdType string, desc string, interpreter string) {
	utils.Log("Inserting Hosts from CommandCapability DB", false)
	stmt, err := database.Connection.Prepare("INSERT INTO `CommandCapability`" +
		"(`cmdParams`, `type`, `name`, `description`, `interpreter`, `displayFields`) VALUES (?, ?, ?, ?, ?, ?);")
	utils.ErrorLog("Couldn't prepare InsertCommand CommandCapability", err, true)

	_, err = stmt.Exec(params, cmdType, name, desc, interpreter, "")
	utils.ErrorLog("Error executing CommandCapability insertHosts", err, true)
	stmt.Close()
}

func UpdateCommandDisplayField(displayFields string, id int) {
	utils.Log("Updating Hosts from CommandCapability DB", false)
	stmt, err := database.Connection.Prepare("UPDATE `CommandCapability` SET `displayFields` = ? WHERE `id` = ?;")
	utils.ErrorLog("Couldn't update CommandCapability", err, true)

	_, err = stmt.Exec(displayFields, id)
	utils.ErrorLog("Results error from CommandCapability", err, true)
	stmt.Close()
}

type Capability struct {
	ID            int
	Params        []string
	Type          string
	Name          string
	Desc          string
	DisplayFields []string
	Interpreter   string
}

func GetAllCapabilities() []Capability {
	utils.Log("Querying capabilities from CommandCapability DB", false)
	stmt, err := database.Connection.Prepare("SELECT `id`, `cmdParams`, `type`, `name`, `description`, `displayFields`, `interpreter` FROM `CommandCapability`")
	utils.ErrorLog("Couldn't select all from CommandCapability GetAllCapabilities", err, true)

	rows, err := stmt.Query()
	utils.ErrorLog("Couldn't recieve rows from CommandCapability GetAllCapabilities", err, true)
	defer rows.Close()

	capabilities := []Capability{}

	for rows.Next() {
		capability := Capability{}
		params := ""
		var paramArr []string

		rows.Scan(&capability.ID, &params, &capability.Type, &capability.Name, &capability.Desc, &capability.DisplayFields, &capability.Interpreter)
		json.Unmarshal([]byte(params), &paramArr)
		capability.Params = paramArr

		capabilities = append(capabilities, capability)
	}

	return capabilities
}
