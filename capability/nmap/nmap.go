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
		}

		// Hostnames
		if len(host.Hostnames) > 0 {
			sysTag := cmdb.EntryTag{
				Label:    "HostName",
				DataType: utils.MAC,
				Values:   []string{},
			}

			for _, name := range host.Hostnames {
				sysTag.Values = append(sysTag.Values, name.Name)
			}

			sysTags = append(sysTags, sysTag)
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

		cmdb.INSERT_ENTRY_Pending(entry)
	}
}
