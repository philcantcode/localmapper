package discovery

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
)

var NMAP_PATH string

func PingScan(ip string) NmapRun {
	utils.Log("Performing NMAP ping scan", true)

	xml, _ := utils.Run("nmap", "-sn", ip, "-oX", "-")
	xmlStruct := ConvertNmapXMLtoStruct(xml)

	utils.PrintLog(xml)

	go insertHosts(xmlStruct)

	return xmlStruct
}

func StealthScan(ip string) NmapRun {
	utils.Log("Performing NMAP Sealth scan (no ping)", true)

	xml, _ := utils.Run("nmap", "-sS", "-Pn", ip, "-oX", "-")
	xmlStruct := ConvertNmapXMLtoStruct(xml)

	go insertHosts(xmlStruct)

	return xmlStruct
}

func OSDetectionScan(ip string) NmapRun {
	utils.Log("Performing NMAP OS Detection scan (no ping)", true)

	xml, _ := utils.Run("nmap", "-O", "-Pn", ip, "-oX", "-")
	xmlStruct := ConvertNmapXMLtoStruct(xml)

	utils.PrintLog(xml)
	utils.PrettyPrint(xmlStruct)

	go insertHosts(xmlStruct)

	return xmlStruct
}

func TestNmap() NmapRun {
	utils.Log("Testing Nmap by scanning 127.0.0.1", false)

	xmlOut, _ := utils.Run("nmap", "-sS", "127.0.0.1", "-oX", "-")
	utils.PrettyPrint(ConvertNmapXMLtoStruct(xmlOut))

	return ConvertNmapXMLtoStruct(xmlOut)
}

func ConvertNmapXMLtoStruct(xmlstr string) NmapRun {
	utils.Log("Converting from XML to an NMAP struct", false)

	var nmapRun NmapRun
	xmlerr := xml.Unmarshal([]byte(xmlstr), &nmapRun)

	if xmlerr != nil {
		utils.ErrorHandle("Error unmarshaling Nmap XML string", xmlerr, false)
	}
	return nmapRun
}

// Check whether a host already exists in DB, update or insert the details
func insertHosts(xml NmapRun) {
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
		utils.ErrorHandle("Couldn't select id from HostTracking", err, true)

		row := stmt.QueryRow(mac, ipv4, ipv6)
		utils.ErrorHandle("Couldn't recieve rows from HostTracking", err, true)

		var id int
		err = row.Scan(&id)

		// Check if no rows returned, in which case, insert the newly found host
		if err == sql.ErrNoRows {
			utils.Log("Inserting Hosts from HostTracking DB", false)
			stmt, err = database.Connection.Prepare("INSERT INTO `HostTracking`" +
				"(`mac`, `ipv4`, `ipv6`, `firstSeen`) VALUES (?, ?, ?, ?);")
			utils.ErrorHandle("Couldn't prepare insertHosts", err, true)

			_, err = stmt.Exec(mac, ipv4, ipv6, time.Now().Unix())
			utils.ErrorHandle("Error executing insertHosts", err, true)
			stmt.Close()
		} else { // If already exists, update it
			utils.Log("Updating Hosts from HostTracking DB", false)
			stmt, err = database.Connection.Prepare("UPDATE `HostTracking` SET `mac` = ?, `ipv4` = ?, `ipv6` = ?, `lastSeen` = ? WHERE `id` = ?;")
			utils.ErrorHandle("Couldn't update HostTracking", err, true)

			_, err = stmt.Exec(mac, ipv4, ipv6, time.Now().Unix(), id)
			utils.ErrorHandle("Results error from HostTracking", err, true)
			stmt.Close()
		}
	}
}

func RegisterNmapCapability() {
	var command string

	fmt.Println("Enter the ")
	fmt.Print("[>] ")
	fmt.Scanf("%s\n", &command)
}
