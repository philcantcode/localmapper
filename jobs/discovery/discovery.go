package network

import (
	"fmt"

	"github.com/philcantcode/localmapper/capability"
	"github.com/philcantcode/localmapper/capability/nmap"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/utils"
)

/* PingSweepVlans performs a sweep of all IP ranges
   from the Vlan database. */
func PingSweepVlans() {
	pingCapability := capability.SELECT_Capability_ByName("Ping Sweep")
	vlans := cmdb.SelectAllVlans()

	// Enumerate each VLan
	for _, vlan := range vlans {
		// Convert the highIP and lowIP to a list of CIDR ranges
		cidr, err := utils.IPv4RangeToCIDRRange(vlan.HighIP, vlan.LowIP)
		utils.ErrorLog(fmt.Sprintf("Couldn't convert IPs to CIDR (%s > %s)", vlan.HighIP, vlan.LowIP), err, true)

		// For each CIDR range, replace the value and run an nmap sweep
		for _, cidrIP := range cidr {
			for i, p := range pingCapability.Command.Params {
				if p.Flag == "" {
					pingCapability.Command.Params[i].Value = cidrIP
				}
			}

			// Log and insert results into DB
			result := nmap.Execute(capability.ParamsToArray(pingCapability.Command.Params))
			nmap.InsertNetworkNmap(result)
			utils.PrintLog(utils.PrettyPrintToStr(result))
		}
	}
}
