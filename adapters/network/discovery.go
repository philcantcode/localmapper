package network

import (
	"fmt"

	"github.com/philcantcode/localmapper/application/database"
	"github.com/philcantcode/localmapper/application/nmap"
	"github.com/philcantcode/localmapper/utils"
)

/* PingSweepVlans performs a sweep of all IP ranges
   from the Vlan database. */
func PingSweepVlans() {
	capability := database.SelectCapability("Ping Sweep")
	vlans := database.SelectAllVlans()

	// Enumerate each VLan
	for _, vlan := range vlans {
		// Convert the highIP and lowIP to a list of CIDR ranges
		cidr, err := utils.IPv4RangeToCIDRRange(vlan.HighIP, vlan.LowIP)
		utils.ErrorLog(fmt.Sprintf("Couldn't convert IPs to CIDR (%s > %s)", vlan.HighIP, vlan.LowIP), err, true)

		// For each CIDR range, replace the value and run an nmap sweep
		for _, cidrIP := range cidr {
			for i, p := range capability.Command.Params {
				if p.Flag == "" {
					capability.Command.Params[i].Value = cidrIP
				}
			}

			// Log and insert results into DB
			result := nmap.RunNmapCommand(capability)
			database.InsertNetworkNmap(result)
			utils.PrintLog(utils.PrettyPrintToStr(result))
		}
	}
}
