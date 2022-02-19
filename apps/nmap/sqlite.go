package nmap

import (
	"database/sql"
	"time"

	"github.com/philcantcode/localmapper/apps/database"
	"github.com/philcantcode/localmapper/utils"
)

// Check whether a host already exists in DB, update or insert the details
func SqliteInsertHosts(xml NmapRun) {
	for i := range xml.Hosts {
		var ipv4 = ""
		var ipv6 = ""
		var mac = ""

		for ii := range xml.Hosts[i].Addresses {

			if xml.Hosts[i].Addresses[ii].AddrType == "ipv4" {
				ipv4 = xml.Hosts[i].Addresses[ii].Addr
			}

			if xml.Hosts[i].Addresses[ii].AddrType == "mac" {
				mac = xml.Hosts[i].Addresses[ii].Addr
			}

			if xml.Hosts[i].Addresses[ii].AddrType == "ipv6" {
				mac = xml.Hosts[i].Addresses[ii].Addr
			}
		}

		// SELECT from the database to see if it already exists
		utils.Log("Querying Hosts from HostTracking DB", false)
		stmt, err := database.Connection.Prepare("SELECT `id` FROM `HostTracking` WHERE" +
			"`mac` = ? OR `ipv4` = ? OR `ipv6` = ? AND `mac` != '' AND ipv4 != '' AND ipv6 != '';")
		utils.ErrorLog("Couldn't select id from HostTracking", err, true)

		row := stmt.QueryRow(mac, ipv4, ipv6)
		utils.ErrorLog("Couldn't recieve rows from HostTracking", err, true)

		var id int
		err = row.Scan(&id)

		// Check if no rows returned, in which case, insert the newly found host
		if err == sql.ErrNoRows {
			utils.Log("Inserting Hosts from HostTracking DB", false)
			stmt, err = database.Connection.Prepare("INSERT INTO `HostTracking`" +
				"(`mac`, `ipv4`, `ipv6`, `firstSeen`) VALUES (?, ?, ?, ?);")
			utils.ErrorLog("Couldn't prepare insertHosts", err, true)

			_, err = stmt.Exec(mac, ipv4, ipv6, time.Now().Unix())
			utils.ErrorLog("Error executing insertHosts", err, true)
			stmt.Close()
		} else { // If already exists, update it
			utils.Log("Updating Hosts from HostTracking DB", false)
			stmt, err = database.Connection.Prepare("UPDATE `HostTracking` SET `mac` = ?, `ipv4` = ?, `ipv6` = ?, `lastSeen` = ? WHERE `id` = ?;")
			utils.ErrorLog("Couldn't update HostTracking", err, true)

			_, err = stmt.Exec(mac, ipv4, ipv6, time.Now().Unix(), id)
			utils.ErrorLog("Results error from HostTracking", err, true)
			stmt.Close()
		}
	}
}
