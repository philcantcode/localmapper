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
	defaultVlanSetup()
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
	prop := Proposition{Type: "local-net-iface", DateTime: local.GetDateTime().DateTime, Desc: "Please choose the IP address for this server.", Predicate: propItem}

	INSERT_Proposition(prop)
}

func recalcualteVlanCIDR() {
	entries := cmdb.SELECT_ENTRY_Inventory(bson.M{"cmdbtype": int32(cmdb.VLAN)}, bson.M{})

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
			entry.SysTags[index] = cmdb.EntryTag{Label: "CIDR", Desc: "CIDR range for this VLAN.", DataType: system.CIDR, Values: cidr}
		} else {
			entry.SysTags = append(entry.SysTags, cmdb.EntryTag{Label: "CIDR", Desc: "CIDR range for this VLAN.", DataType: system.CIDR, Values: cidr})
		}

		cmdb.UPDATE_ENTRY_Inventory(entry)
	}
}

func defaultVlanSetup() {
	vlan1 := cmdb.SELECT_ENTRY_Inventory(bson.M{"label": "Private Range 1", "desc": "Default VLAN"}, bson.M{})

	if len(vlan1) == 0 {
		lowIP := cmdb.EntryTag{Label: "LowIP", DataType: system.IP_RANGE_LOW, Values: []string{"10.0.0.0"}}
		highIP := cmdb.EntryTag{Label: "HighIP", DataType: system.IP_RANGE_HIGH, Values: []string{"10.255.255.255"}}
		sysDefault := cmdb.EntryTag{Label: "SysDefault", DataType: system.BOOL, Values: []string{"1"}}

		newVlan := cmdb.Entry{Label: "Private Range 1", Desc: "Default VLAN", CMDBType: cmdb.VLAN, OSILayer: 2, DateSeen: []string{utils.Now()}, SysTags: []cmdb.EntryTag{lowIP, highIP, sysDefault}}
		cmdb.INSERT_ENTRY_Inventory(newVlan)
	}

	vlan2 := cmdb.SELECT_ENTRY_Inventory(bson.M{"label": "Private Range 2", "desc": "Default VLAN"}, bson.M{})

	if len(vlan2) == 0 {
		lowIP := cmdb.EntryTag{Label: "LowIP", DataType: system.IP_RANGE_LOW, Values: []string{"172.16.0.0"}}
		highIP := cmdb.EntryTag{Label: "HighIP", DataType: system.IP_RANGE_HIGH, Values: []string{"172.31.255.255"}}
		sysDefault := cmdb.EntryTag{Label: "SysDefault", DataType: system.BOOL, Values: []string{"1"}}

		newVlan := cmdb.Entry{Label: "Private Range 2", Desc: "Default VLAN", CMDBType: cmdb.VLAN, OSILayer: 2, DateSeen: []string{utils.Now()}, SysTags: []cmdb.EntryTag{lowIP, highIP, sysDefault}}
		cmdb.INSERT_ENTRY_Inventory(newVlan)
	}

	vlan3 := cmdb.SELECT_ENTRY_Inventory(bson.M{"label": "Private Range 3", "desc": "Default VLAN"}, bson.M{})

	if len(vlan3) == 0 {
		lowIP := cmdb.EntryTag{Label: "LowIP", DataType: system.IP_RANGE_LOW, Values: []string{"192.168.0.0"}}
		highIP := cmdb.EntryTag{Label: "HighIP", DataType: system.IP_RANGE_HIGH, Values: []string{"192.168.255.255"}}
		sysDefault := cmdb.EntryTag{Label: "SysDefault", DataType: system.BOOL, Values: []string{"1"}}

		newVlan := cmdb.Entry{Label: "Private Range 3", Desc: "Default VLAN", CMDBType: cmdb.VLAN, OSILayer: 2, DateSeen: []string{utils.Now()}, SysTags: []cmdb.EntryTag{lowIP, highIP, sysDefault}}
		cmdb.INSERT_ENTRY_Inventory(newVlan)
	}

	vlan4 := cmdb.SELECT_ENTRY_Inventory(bson.M{"label": "Test Home", "desc": "Test VLAN"}, bson.M{})

	if len(vlan4) == 0 {
		lowIP := cmdb.EntryTag{Label: "LowIP", DataType: system.IP_RANGE_LOW, Values: []string{"192.168.1.0"}}
		highIP := cmdb.EntryTag{Label: "HighIP", DataType: system.IP_RANGE_HIGH, Values: []string{"192.168.1.255"}}
		sysDefault := cmdb.EntryTag{Label: "SysDefault", DataType: system.BOOL, Values: []string{"1"}}

		newVlan := cmdb.Entry{Label: "Test Home", Desc: "Test VLAN", CMDBType: cmdb.VLAN, OSILayer: 2, DateSeen: []string{utils.Now()}, SysTags: []cmdb.EntryTag{lowIP, highIP, sysDefault}}
		cmdb.INSERT_ENTRY_Inventory(newVlan)
	}
}
