package propositions

import (
	"github.com/philcantcode/localmapper/capabilities/local"
	"github.com/philcantcode/localmapper/core"
	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
)

func SetupJobs() {
	setupSelfIdentity()
}

// setupSelfIdentity initialises IPs and local variables
func setupSelfIdentity() {
	cmdbs := database.SelectAllCMDB()

	// Check if server is already in the database
	for _, cmdb := range cmdbs {
		ident, found := cmdb.InfoTags["identity"]

		if found && ident == "local" {
			utils.Log("Identity local already found in CMDB", true)
			return
		}
	}

	// Delete old propositions from previous reboots
	for _, proposition := range database.SelectAllPropositions() {
		if proposition.Type == "local-net-iface" {
			database.SetPropositionStatusByID(proposition.ID, 2) // 2 = disabled
		}
	}

	// Create the new proposition to guess the local IP
	optionIPs := []string{}

	for _, ip := range local.GetNetworkAdapters() {
		optionIPs = append(optionIPs, ip.IP)
	}

	propItem := core.PropositionItem{Name: "Server IP", Value: local.GetDefaultIPGateway().DefaultIP, DataType: core.IP, Options: optionIPs}
	prop := core.Proposition{Type: "local-net-iface", Date: local.GetDateTime().DateTime, Description: "Please choose the IP address for this server.", Correction: core.PropositionItem{}, Proposition: propItem}

	database.InsertProposition(prop)
}
