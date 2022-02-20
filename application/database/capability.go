package database

import (
	"encoding/json"

	"github.com/philcantcode/localmapper/adapters/blueprint"
	"github.com/philcantcode/localmapper/utils"
)

func InsertCapability(capability blueprint.Capability) {
	utils.Log("Inserting Hosts from CommandCapability DB", false)
	stmt, err := connection.Prepare("INSERT INTO `Capabilities`" +
		"(`cmdParams`, `type`, `name`, `description`, `displayFields`) VALUES (?, ?, ?, ?, ?);")
	utils.ErrorLog("Couldn't prepare InsertCommand CommandCapability", err, true)

	_, err = stmt.Exec(capability.Params, capability.Type, capability.Name, capability.Desc, "")
	utils.ErrorLog("Error executing CommandCapability insertHosts", err, true)
	stmt.Close()
}

func UpdateCapability(capability blueprint.Capability) {
	utils.Log("Updating Capabilities database", false)
	stmt, err := connection.Prepare("UPDATE `Capabilities` SET `name` = ?, `type` = ?, `params` = ?, `description` = ?, `displayFields` = ? WHERE `id` = ?;")
	utils.ErrorFatal("Couldn't update Capabilities database", err)

	_, err = stmt.Exec(capability.Name, capability.Type, capability.Params, capability.Desc, capability.DisplayFields, capability.ID)
	utils.ErrorFatal("Results error from UpdateCapability", err)
	stmt.Close()
}

func SelectAllCapabilities() []blueprint.Capability {
	utils.Log("Querying capabilities from Capabilities DB", false)
	stmt, err := connection.Prepare("SELECT `id`, `params`, `type`, `name`, `description`, `displayFields` FROM `Capabilities`")
	utils.ErrorLog("Couldn't select all from Capabilities GetAllCapabilities", err, true)

	rows, err := stmt.Query()
	utils.ErrorLog("Couldn't recieve rows from CommandCapability GetAllCapabilities", err, true)
	defer rows.Close()

	capabilities := []blueprint.Capability{}

	for rows.Next() {
		capability := blueprint.Capability{}
		params := ""
		var paramArr []string

		rows.Scan(&capability.ID, &params, &capability.Type, &capability.Name, &capability.Desc, &capability.DisplayFields)
		json.Unmarshal([]byte(params), &paramArr)
		capability.Params = paramArr

		capabilities = append(capabilities, capability)
	}

	return capabilities
}
