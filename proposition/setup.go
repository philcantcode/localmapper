package proposition

import (
	"github.com/philcantcode/localmapper/capability/local"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func SetupJobs() {
	setupSelfIdentity()
	calculateVlanCIDR()
}

// setupSelfIdentity initialises IPs and local variables
func setupSelfIdentity() {
	entries := cmdb.SELECT_ENTRY(bson.M{}, bson.M{})
	localInCMDB := false

	// Check if server is already in the database
	for _, entry := range entries {
		ident, found := cmdb.FindSysTag("identity", entry)

		if found && utils.ArrayContains("local", ident.Values) {
			utils.Log("Identity local already found in CMDB", true)
			return
		}
	}

	props := SELECT_Propositions(bson.M{"type": "local-net-iface"}, bson.M{})

	if localInCMDB && len(props) > 0 {
		utils.Log("Proposition for self identity already exists", true)
		return
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

func calculateVlanCIDR() {
	entries := cmdb.SELECT_ENTRY(bson.M{"cmdbtype": int32(cmdb.VLAN)}, bson.M{})

	for _, entry := range entries {
		// Check CMDB entry is of type VLAN
		if entry.CMDBType != cmdb.VLAN {
			continue
		}

		lowIP, lowFound := cmdb.FindUsrTag("LowIP", entry)
		highIP, highFound := cmdb.FindUsrTag("HighIP", entry)

		// Check that both of the user tags for the IPs are set
		if !lowFound && !highFound {
			continue
		}

		cidr, err := utils.IPv4RangeToCIDRRange(lowIP.Values[0], highIP.Values[0])
		utils.ErrorLog("Couldn't generate CIDR for: "+entry.Label, err, true)

		// Remove old CMDB tags so new one can be calcualted
		entry.SysTags = cmdb.RemoveTag(entry.SysTags, "CMDB")

		entry.SysTags = append(entry.SysTags, cmdb.EntryTag{Label: "CIDR", Desc: "CIDR range for this VLAN.", DataType: utils.CIDR, Values: cidr})
		cmdb.UPDATE_ENTRY(entry)
	}
}
