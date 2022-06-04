package proposition

import (
	"github.com/philcantcode/localmapper/capability/local"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/utils"
)

func SetupJobs() {
	setupSelfIdentity()
}

// setupSelfIdentity initialises IPs and local variables
func setupSelfIdentity() {
	cmdbs := cmdb.SelectAllCMDB()

	// Check if server is already in the database
	for _, cmdb := range cmdbs {
		ident, found := cmdb.InfoTags["identity"]

		if found && ident == "local" {
			utils.Log("Identity local already found in CMDB", true)
			return
		}
	}

	// Delete old propositions from previous reboots
	for _, proposition := range SELECT_Propositions_All() {
		if proposition.Type == "local-net-iface" {
			UPDATE_Proposition_Status_ByID(proposition.ID, 2) // 2 = disabled
		}
	}

	// Create the new proposition to guess the local IP
	optionIPs := []string{}

	for _, ip := range local.GetNetworkAdapters() {
		optionIPs = append(optionIPs, ip.IP)
	}

	propItem := PropositionItem{Name: "Server IP", Value: local.GetDefaultIPGateway().DefaultIP, DataType: utils.IP, Options: optionIPs}
	prop := Proposition{Type: "local-net-iface", Date: local.GetDateTime().DateTime, Description: "Please choose the IP address for this server.", Correction: PropositionItem{}, Proposition: propItem}

	INSERT_Proposition(prop)
}
