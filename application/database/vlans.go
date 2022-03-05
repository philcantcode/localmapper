package database

import (
	"github.com/philcantcode/localmapper/adapters/definitions"
	"github.com/philcantcode/localmapper/utils"
)

func SelectAllVlans() []definitions.Vlan {
	utils.Log("SelectAllVlans from Vlans Db (sqlite)", false)
	stmt, err := connection.Prepare("SELECT `id`, `name`, `description`, `highIP`, `lowIP` FROM `Vlans`")
	utils.ErrorLog("Couldn't select all from Vlans", err, true)

	rows, err := stmt.Query()
	utils.ErrorLog("Couldn't recieve rows from SelectAllVlans", err, true)
	defer rows.Close()

	vlans := []definitions.Vlan{}

	for rows.Next() {
		vlan := definitions.Vlan{}

		rows.Scan(&vlan.ID, &vlan.Name, &vlan.Description, &vlan.HighIP, &vlan.LowIP)

		vlans = append(vlans, vlan)
	}

	return vlans
}
