package network

import (
	"fmt"
	"net"

	"github.com/philcantcode/localmapper/capabilities/local"
	"github.com/philcantcode/localmapper/capabilities/nmap"
	"github.com/philcantcode/localmapper/database"
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

/* GenerateListOfGatewaysFromNetworkAdapters calculates the first IP on every
   network adapter attached. */
func GenerateListOfGatewaysFromNetworkAdapters() map[string]string {
	netAdapters := local.GetNetworkAdapters()

	// {Adapter Name : IP Address}
	for key, addr := range netAdapters {
		ip := []byte(net.ParseIP(addr).To4())
		gateway := []byte{ip[0], ip[1], ip[2], 1}

		fmt.Printf("%d -> %d\n", ip, gateway)
		netAdapters[key] = string(net.IPv4(ip[0], ip[1], ip[2], 1).String())
		utils.Log(fmt.Sprintf("Calculating gateway: %d -> %d", ip, gateway), false)
	}

	return netAdapters
}
