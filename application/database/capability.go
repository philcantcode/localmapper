package database

import (
	"encoding/json"

	"github.com/philcantcode/localmapper/adapters/definitions"
	"github.com/philcantcode/localmapper/utils"
)

func InsertCapability(capability definitions.Capability) {
	utils.Log("Inserting Hosts from CommandCapability DB", false)
	stmt, err := connection.Prepare("INSERT INTO `Capabilities`" +
		"(`name`, `type`, `command`, `description`, `displayFields`) VALUES (?, ?, ?, ?, ?);")
	utils.ErrorLog("Couldn't prepare InsertCommand CommandCapability", err, true)

	command, err := json.Marshal(capability.Command)
	utils.ErrorFatal("Couldn't Marshall Capability Command", err)

	_, err = stmt.Exec(capability.Name, capability.Type, string(command), capability.Desc, capability.DisplayFields)
	utils.ErrorLog("Error executing CommandCapability insertHosts", err, true)
	stmt.Close()
}

func UpdateCapability(capability definitions.Capability) {
	utils.Log("Updating Capabilities database", false)
	stmt, err := connection.Prepare("UPDATE `Capabilities` SET `name` = ?, `type` = ?, `command` = ?, `description` = ?, `displayFields` = ? WHERE `id` = ?;")
	utils.ErrorFatal("Couldn't update Capabilities database", err)

	command, err := json.Marshal(capability.Command)
	utils.ErrorFatal("Couldn't Marshall Capability Command", err)

	_, err = stmt.Exec(capability.Name, capability.Type, string(command), capability.Desc, capability.DisplayFields, capability.ID)
	utils.ErrorFatal("Results error from UpdateCapability", err)
	stmt.Close()
}

func SelectAllCapabilities() []definitions.Capability {
	utils.Log("Querying capabilities from Capabilities DB", false)
	stmt, err := connection.Prepare("SELECT `id`, `command`, `type`, `name`, `description`, `displayFields` FROM `Capabilities`")
	utils.ErrorLog("Couldn't select all from Capabilities GetAllCapabilities", err, true)

	rows, err := stmt.Query()
	utils.ErrorLog("Couldn't recieve rows from SelectAllCapabilities", err, true)
	defer rows.Close()

	capabilities := []definitions.Capability{}

	for rows.Next() {
		capability := definitions.Capability{}
		command := ""

		rows.Scan(&capability.ID, &command, &capability.Type, &capability.Name, &capability.Desc, &capability.DisplayFields)
		json.Unmarshal([]byte(command), &capability.Command)

		capabilities = append(capabilities, capability)
	}

	return capabilities
}

func SelectCapability(name string) definitions.Capability {
	capabilities := SelectAllCapabilities()

	for _, k := range capabilities {
		if k.Name == name {
			return k
		}
	}

	utils.ErrorForceFatal("Could not SelectCapability for: " + name)
	return definitions.Capability{}
}
