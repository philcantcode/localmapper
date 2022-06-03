package database

import (
	"github.com/philcantcode/localmapper/core"
	"github.com/philcantcode/localmapper/utils"
)

// Descending Order
func SelectAllVlans() []core.Vlan {
	utils.Log("SelectAllVlans from Vlans Db (sqlite)", false)
	stmt, err := connection.Prepare("SELECT `id`, `name`, `description`, `highIP`, `lowIP`, `tags` FROM `Vlans` ORDER BY `id` DESC")
	utils.ErrorLog("Couldn't select all from Vlans", err, true)

	rows, err := stmt.Query()
	utils.ErrorLog("Couldn't recieve rows from SelectAllVlans", err, true)
	defer rows.Close()

	vlans := []core.Vlan{}

	for rows.Next() {
		vlan := core.Vlan{}

		rows.Scan(&vlan.ID, &vlan.Name, &vlan.Description, &vlan.HighIP, &vlan.LowIP, &vlan.Tags)

		vlans = append(vlans, vlan)
	}

	return vlans
}
