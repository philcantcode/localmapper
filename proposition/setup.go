package proposition

import (
	"github.com/philcantcode/localmapper/capability/local"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func FirstTimeSetup() {
	setupSelfIdentity()
	recalcualteVlanCIDR()
}

// setupSelfIdentity initialises IPs and local variables
func setupSelfIdentity() {
	entries := cmdb.SELECT_ENTRY_Inventory(bson.M{}, bson.M{})
	localInCMDB := false
	propAlreadyExists := false

	// Check if server is already in the database
	for _, entry := range entries {
		ident, found, _ := cmdb.FindSysTag("Identity", entry)

		if found && utils.ArrayContains("local", ident.Values) {
			system.Log("Identity local already found in CMDB", true)
			localInCMDB = true
		}
	}

	props := SELECT_Propositions(bson.M{"type": "local-net-iface"}, bson.M{})

	if len(props) > 0 {
		system.Log("Proposition for self identity already exists", true)
		propAlreadyExists = true
	}

	if localInCMDB {
		return
	}

	if localInCMDB && propAlreadyExists {
		return
	}

	for _, prop := range props {
		prop.Status = 2
		UPDATE_Proposition(prop)
	}

	// Create the new proposition to guess the local IP
	optionIPs := []string{}

	for _, ip := range local.GetNetworkAdapters() {
		optionIPs = append(optionIPs, ip.IP)
	}

	propItem := Predicate{Label: "Server IP", Value: local.GetDefaultIPGateway().DefaultIP, DataType: system.IP, Options: optionIPs}
	prop := Proposition{Type: "local-net-iface", DateTime: local.GetDateTime().DateTime, Description: "Please choose the IP address for this server.", Predicate: propItem}

	INSERT_Proposition(prop)
}

func recalcualteVlanCIDR() {
	entries := cmdb.SELECT_ENTRY_Inventory(bson.M{"cmdbtype": cmdb.VLAN}, bson.M{})

	for _, entry := range entries {
		// Check CMDB entry is of type VLAN
		if entry.CMDBType != cmdb.VLAN {
			continue
		}

		lowIP, lowFound, _ := cmdb.FindSysTag("LowIP", entry)
		highIP, highFound, _ := cmdb.FindSysTag("HighIP", entry)

		// Check that both of the user tags for the IPs are set
		if !lowFound && !highFound {
			continue
		}

		cidr, err := utils.IPv4RangeToCIDRRange(lowIP.Values[0], highIP.Values[0])
		system.Error("Couldn't generate CIDR for: "+entry.Label, err)

		// Remove old CMDB tags so new one can be calcualted
		_, found, index := cmdb.FindSysTag("CIDR", entry)

		if found {
			entry.SysTags[index] = cmdb.EntryTag{Label: "CIDR", Description: "CIDR range for this VLAN.", DataType: system.CIDR, Values: cidr}
		} else {
			entry.SysTags = append(entry.SysTags, cmdb.EntryTag{Label: "CIDR", Description: "CIDR range for this VLAN.", DataType: system.CIDR, Values: cidr})
		}

		cmdb.UPDATE_ENTRY_Inventory(entry)
	}
}
