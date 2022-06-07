package nmap

import (
	"encoding/xml"
	"fmt"
	"os/exec"

	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/utils"
)

// Execute takes a list of parameters to execute against NMAP
func Execute(params []string) NmapRun {
	utils.Log(fmt.Sprintf("Attempting to run Nmap Command: %s > %v", "nmap", params), true)
	resultByte, err := exec.Command("nmap", params...).CombinedOutput()
	utils.ErrorFatal(fmt.Sprintf("Error returned running a command: %s > %v", "nmap", params), err)

	utils.Log("Converting from []byte to NmapRun struct", false)

	var nmapRun NmapRun
	err = xml.Unmarshal(resultByte, &nmapRun)
	utils.ErrorLog("Couldn't unmarshal result from Nmap console", err, true)

	interpret(nmapRun)

	return nmapRun
}

func interpret(nmapRun NmapRun) {
	// For each host
	for _, host := range nmapRun.Hosts {
		sysTags := []cmdb.EntryTag{}
		vendorTags := cmdb.EntryTag{
			Label:    "Vendor",
			DataType: utils.IP,
			Values:   []string{},
		}

		for _, address := range host.Addresses {
			if address.AddrType == "ipv4" {
				sysTags = append(sysTags, cmdb.EntryTag{
					Label:    "IP",
					DataType: utils.IP,
					Values:   []string{address.Addr},
				})
			}

			if address.AddrType == "mac" {
				sysTags = append(sysTags, cmdb.EntryTag{
					Label:    "MAC",
					DataType: utils.MAC,
					Values:   []string{address.Addr},
				})
			}

			if address.Vendor != "" {
				vendorTags.Values = append(vendorTags.Values, address.Vendor)
			}
		}

		if len(vendorTags.Values) > 0 {
			sysTags = append(sysTags, vendorTags)
		}

		// Hostnames
		if len(host.Hostnames) > 0 {
			hostNameTag := cmdb.EntryTag{
				Label:    "HostName",
				DataType: utils.MAC,
				Values:   []string{},
			}

			for _, name := range host.Hostnames {
				hostNameTag.Values = append(hostNameTag.Values, name.Name)
			}

			sysTags = append(sysTags, hostNameTag)
		}

		entry := cmdb.Entry{
			Label:    "Nmap Discovered Device",
			Desc:     "This device was discovered during an Nmap scan: " + nmapRun.Args,
			OSILayer: 0,
			CMDBType: cmdb.PENDING,
			DateSeen: []string{nmapRun.StartStr},
			SysTags:  sysTags,
		}

		tag, exists, _ := cmdb.FindSysTag("HostName", entry)

		if exists {
			entry.Label = tag.Values[0]
		}

		// Insert to pending or update both DBs
		if cmdb.EntryExists_ByIP(entry) {
			newInventoy := cmdb.UpdateInventoryEntries_ByIP(entry)
			newPending := cmdb.UpdatePendingEntries_ByIP(entry)

			if !newInventoy && !newPending {
				utils.FatalAlert("Couldn't update inventory or pending in nmap")
			}
		} else {
			cmdb.INSERT_ENTRY_Pending(entry)
		}
	}
}
