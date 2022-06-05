package proposition

import (
	"github.com/philcantcode/localmapper/capability/local"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func SetupJobs() {
	setupSelfIdentity()
}

// setupSelfIdentity initialises IPs and local variables
func setupSelfIdentity() {
	cmdbs := cmdb.SELECT_CMDBItem_All()

	// Check if server is already in the database
	for _, cmdb := range cmdbs {
		ident, found := cmdb.InfoTags["identity"]

		if found && ident == "local" {
			utils.Log("Identity local already found in CMDB", true)
			return
		}
	}

	// Delete old propositions from previous reboots
	for _, proposition := range SELECT_Propositions(bson.M{}, bson.M{}) {
		if proposition.Type == "local-net-iface" {
			proposition.Status = 2
			UPDATE_Proposition(proposition) // 2 = disabled
		}
	}

	// Create the new proposition to guess the local IP
	optionIPs := []string{}

	for _, ip := range local.GetNetworkAdapters() {
		optionIPs = append(optionIPs, ip.IP)
	}

	propItem := Predicate{Label: "Server IP", Value: local.GetDefaultIPGateway().DefaultIP, DataType: utils.IP, Options: optionIPs}
	prop := Proposition{Type: "local-net-iface", DateTime: local.GetDateTime().DateTime, Desc: "Please choose the IP address for this server.", Predicate: propItem}

	INSERT_Proposition(prop)
}
