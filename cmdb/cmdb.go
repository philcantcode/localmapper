package cmdb

import (
	"fmt"
	"strconv"
	"time"

	"github.com/philcantcode/localmapper/local"
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

/*
	EntryExistsSomewhere returns true if an entry exists either
	Inventory OR Pending
*/
func EntryExists_ByIP(entry Entity) bool {
	tag, exists, _ := entry.FindSysTag("IP")
	result := []Entity{}

	if !exists {
		return false
	}

	ipFilter := bson.M{
		"systags.label":  "IP",
		"systags.values": tag.Values[len(tag.Values)-1],
	}

	result = append(result, SELECT_ENTRY_Inventory(ipFilter, bson.M{})...)
	result = append(result, SELECT_ENTRY_Pending(ipFilter, bson.M{})...)

	// Return true if len > 0, else false
	return len(result) != 0
}

/*
	Updates both database entries by IP.

	Finds the old entry result[0] and then updates the
	values to the new entry.
*/
func updateEntriesTags_ByIP(entry Entity) bool {
	tag, exists, _ := entry.FindSysTag("IP")

	if !exists {
		return false
	}

	ipFilter := bson.M{
		"systags.label":  "IP",
		"systags.values": tag.Values[len(tag.Values)-1],
	}

	results := SELECT_Entities_Joined(ipFilter, bson.M{})

	if len(results) == 0 {
		system.Log(fmt.Sprintf("No match for (inventory): %s", tag.Values[len(tag.Values)-1]), false)
		return false
	}

	if len(results) > 1 { // Too many results returned, database corrupt
		system.Warning(
			fmt.Sprintf(
				"While executing UpdateInventoryEntries the number of matched results > 1\n\nEntry: %+v\n\nMatched Cases: %+v", entry, results), false)

		InitIPConflict(SELECT_ENTRY_Pending(ipFilter, bson.M{})[0], SELECT_ENTRY_Inventory(ipFilter, bson.M{})[0])
		return false
	}

	system.Log(fmt.Sprintf("Match (Inventory): len: %d, IP: %+v", len(results), results), false)

	// Parse SysTags and join them
	for _, newTag := range entry.SysTags {
		_, found, i := results[0].FindSysTag(newTag.Label)

		if found {
			results[0].SysTags[i].Values = joinTagGroups(newTag.Label, results[0].SysTags[i].Values, newTag.Values)
		} else {
			results[0].SysTags = append(results[0].SysTags, newTag)
		}
	}

	// Parse SysTags and join them
	for _, newTag := range entry.UsrTags {
		_, found, i := results[0].FindUsrTag(newTag.Label)

		if found {
			results[0].UsrTags[i].Values = joinTagGroups(newTag.Label, results[0].UsrTags[i].Values, newTag.Values)
		} else {
			results[0].UsrTags = append(results[0].UsrTags, newTag)
		}
	}

	results[0].DateSeen = append(results[0].DateSeen, entry.DateSeen...)

	system.Log(fmt.Sprintf("Compartive update made: %v", results[0].ID), false)
	results[0].UPDATE_ENTRY_Inventory()

	// Only update the metadata for the pending entry
	results[0].Label = entry.Label
	results[0].Description = entry.Description
	results[0].CMDBType = entry.CMDBType
	results[0].OSILayer = entry.OSILayer

	results[0].UPDATE_ENTRY_Pending()

	return true
}

/*

 */
func joinTagGroups(label string, oldTags []string, newTags []string) []string {
	if len(newTags) == 0 {
		return oldTags
	}

	if len(oldTags) == 0 {
		return newTags
	}

	// Update old tags with new tags
	switch label {
	case "Ports", "Services": // Merge uniques
		for _, val := range newTags {
			if !utils.ArrayContains(val, oldTags) {
				oldTags = append(oldTags, val)
			}
		}
	case "IP", "IP6", "MAC", "MAC6": // Append if last element not same
		for _, val := range newTags {
			if oldTags[len(oldTags)-1] != val {
				oldTags = append(oldTags, val)
			}
		}
	default: // Overwrite
		oldTags = newTags
	}

	return oldTags
}

type IdentityConfidence struct {
	IP       int
	MAC      int
	HostName int
	Vendor   int
	OS       int
	DateSeen int

	Total   int
	Average int
}

/*
	CalcIdentityConfidenceScore returns an structure of confidences for various
	values known about the device.
*/
func CalcIdentityConfidenceScore(entry Entity) IdentityConfidence {
	result := IdentityConfidence{IP: 0, MAC: 0, HostName: 0, Vendor: 0, OS: 0, DateSeen: 0, Average: 0}

	// IP4 weighted 80%, divided by the number of past changes
	ipTag, found, _ := entry.FindSysTag("IP")

	if found {
		result.IP += (80 / len(ipTag.Values))
	}

	// IP6 weigthed 20%, divided by the number of past changes
	ip6Tag, found, _ := entry.FindSysTag("IP6")

	if found {
		result.IP += (20 / len(ip6Tag.Values))
	}

	// MAC weighted 80%, divided by the number of past changes
	macTag, found, _ := entry.FindSysTag("MAC")

	if found {
		result.MAC += (80 / len(macTag.Values))
	}

	// MAC6 weigthed 20%, divided by the number of past changes
	mac6Tag, found, _ := entry.FindSysTag("MAC6")

	if found {
		result.MAC += (20 / len(mac6Tag.Values))
	}

	// Host Name weigthed 100%, divided by the number of past changes
	hostNameTag, found, _ := entry.FindSysTag("HostName")

	if found {
		result.HostName = (100 / len(hostNameTag.Values))
	}

	_, hasMACVendor, _ := entry.FindSysTag("MACVendor")
	_, hasOSVendor, _ := entry.FindSysTag("OSVendor")
	vendorACC, hasVendorACC, _ := entry.FindSysTag("OSAccuracy")

	if hasVendorACC && hasOSVendor {
		vAccInt, err := strconv.Atoi(vendorACC.Values[0])
		system.Fatal("Couldn't convert CMDBType to int", err)

		result.Vendor += (vAccInt / 2)
	}

	if hasMACVendor {
		result.Vendor += 50
	}

	_, found, _ = entry.FindSysTag("OS")

	if found {
		result.OS = 100
	}

	// Assign values based on the number of dates seen
	size := len(entry.DateSeen)
	if size == 0 {
		result.DateSeen = 0
	} else if size > 0 && size <= 5 {
		result.DateSeen = 33
	} else if size > 5 && size <= 10 {
		result.DateSeen = 50
	} else if size > 10 && size <= 50 {
		result.DateSeen = 75
	} else if size > 50 {
		result.DateSeen = 100
	}

	// Calcualte an average
	result.Total = result.IP + result.MAC + result.OS + result.Vendor + result.DateSeen + result.HostName
	result.Average = result.Total / 6

	return result
}

func CalcTimeGraph(entry Entity) TimeGraph {
	graph := TimeGraph{Keys: []string{}, Values: []int{}}

	windowMinute, err := strconv.Atoi(system.GetConfig("date-seen-graph-mins-val"))
	system.Fatal("Couldn't convert date-seen-graph-mins-val settings to int", err)

	window := time.Duration(int64(time.Minute) * int64(windowMinute))
	nowDT, _ := time.Parse(utils.DTF_DateTime, entry.DateSeen[0])
	endDT, _ := time.Parse(utils.DTF_DateTime, entry.DateSeen[len(entry.DateSeen)-1])
	processedBlocks := 1
	entriesInBlock := 0
	timeBlockCount := 0

	// Push the first value of DateTime to get things started
	graph.Keys = append(graph.Keys, nowDT.Format(utils.DTF_DateTime))

	for nowDT.Before(endDT) {
		nowDT = nowDT.Add(window) // Add 1 duration step to it

		for _, dt := range entry.DateSeen[processedBlocks:] {
			subDT, _ := time.Parse(utils.DTF_DateTime, dt)

			if subDT.Before(nowDT) {
				entriesInBlock++
			} else {
				break
			}

			processedBlocks++
		}

		graph.Keys = append(graph.Keys, nowDT.Format(utils.DTF_DateTime))
		graph.Values = append(graph.Values, entriesInBlock)

		entriesInBlock = 0
		timeBlockCount++
	}

	return graph
}

func Init() {

	InitLocalIdentityProp()
	updateSelfIdentity()

	if len(SELECT_ENTRY_Inventory(bson.M{}, bson.M{})) > 0 {
		return
	}

	vlan1 := SELECT_ENTRY_Inventory(bson.M{"label": "Private Range 1", "desc": "Default VLAN"}, bson.M{})

	if len(vlan1) == 0 {
		lowIP := EntityTag{Label: "LowIP", DataType: system.DataType_IP_RANGE_LOW, Values: []string{"10.0.0.0"}}
		highIP := EntityTag{Label: "HighIP", DataType: system.DataType_IP_RANGE_HIGH, Values: []string{"10.255.255.255"}}
		sysDefault := EntityTag{Label: "SysDefault", DataType: system.DataType_BOOL, Values: []string{"1"}}

		newVlan := Entity{Label: "Private Range 1", Description: "Default VLAN", CMDBType: VLAN, OSILayer: 2, DateSeen: []string{utils.Now()}, SysTags: []EntityTag{lowIP, highIP, sysDefault}}
		newVlan.InsertInventory()
	}

	vlan2 := SELECT_ENTRY_Inventory(bson.M{"label": "Private Range 2", "desc": "Default VLAN"}, bson.M{})

	if len(vlan2) == 0 {
		lowIP := EntityTag{Label: "LowIP", DataType: system.DataType_IP_RANGE_LOW, Values: []string{"172.16.0.0"}}
		highIP := EntityTag{Label: "HighIP", DataType: system.DataType_IP_RANGE_HIGH, Values: []string{"172.31.255.255"}}
		sysDefault := EntityTag{Label: "SysDefault", DataType: system.DataType_BOOL, Values: []string{"1"}}

		newVlan := Entity{Label: "Private Range 2", Description: "Default VLAN", CMDBType: VLAN, OSILayer: 2, DateSeen: []string{utils.Now()}, SysTags: []EntityTag{lowIP, highIP, sysDefault}}
		newVlan.InsertInventory()
	}

	vlan3 := SELECT_ENTRY_Inventory(bson.M{"label": "Private Range 3", "desc": "Default VLAN"}, bson.M{})

	if len(vlan3) == 0 {
		lowIP := EntityTag{Label: "LowIP", DataType: system.DataType_IP_RANGE_LOW, Values: []string{"192.168.0.0"}}
		highIP := EntityTag{Label: "HighIP", DataType: system.DataType_IP_RANGE_HIGH, Values: []string{"192.168.255.255"}}
		sysDefault := EntityTag{Label: "SysDefault", DataType: system.DataType_BOOL, Values: []string{"1"}}

		newVlan := Entity{Label: "Private Range 3", Description: "Default VLAN", CMDBType: VLAN, OSILayer: 2, DateSeen: []string{utils.Now()}, SysTags: []EntityTag{lowIP, highIP, sysDefault}}
		newVlan.InsertInventory()
	}

	vlan4 := SELECT_ENTRY_Inventory(bson.M{"label": "Test Home", "desc": "Test VLAN"}, bson.M{})

	if len(vlan4) == 0 {
		lowIP := EntityTag{Label: "LowIP", DataType: system.DataType_IP_RANGE_LOW, Values: []string{"192.168.1.0"}}
		highIP := EntityTag{Label: "HighIP", DataType: system.DataType_IP_RANGE_HIGH, Values: []string{"192.168.1.255"}}
		sysDefault := EntityTag{Label: "SysDefault", DataType: system.DataType_BOOL, Values: []string{"1"}}

		newVlan := Entity{Label: "Test Home", Description: "Test VLAN", CMDBType: VLAN, OSILayer: 2, DateSeen: []string{utils.Now()}, SysTags: []EntityTag{lowIP, highIP, sysDefault}}
		newVlan.InsertInventory()
	}

	vlan5 := SELECT_ENTRY_Inventory(bson.M{"label": "Olivers Home", "desc": "Test VLAN"}, bson.M{})

	if len(vlan5) == 0 {
		lowIP := EntityTag{Label: "LowIP", DataType: system.DataType_IP_RANGE_LOW, Values: []string{"192.168.0.0"}}
		highIP := EntityTag{Label: "HighIP", DataType: system.DataType_IP_RANGE_HIGH, Values: []string{"192.168.0.255"}}
		sysDefault := EntityTag{Label: "SysDefault", DataType: system.DataType_BOOL, Values: []string{"1"}}

		newVlan := Entity{Label: "Olivers Home", Description: "Test VLAN", CMDBType: VLAN, OSILayer: 2, DateSeen: []string{utils.Now()}, SysTags: []EntityTag{lowIP, highIP, sysDefault}}
		newVlan.InsertInventory()
	}

	recalcualteVlanCIDR()
}

/*
	UpdateOrInsert either updates the entity by IP if found
	or inserts a new entity. Causes DATE_SEEN to be updated
*/
func (entry Entity) UpdateOrInsert() {
	// Insert to pending or update both DBs
	if EntryExists_ByIP(entry) {
		entryUpdateSuccess := updateEntriesTags_ByIP(entry)

		if !entryUpdateSuccess {
			system.Warning("Couldn't update inventory or pending in CMDB", false)
		}
	} else {
		entry.InsertPending()
	}
}

func recalcualteVlanCIDR() {
	entries := SELECT_ENTRY_Inventory(bson.M{"cmdbtype": VLAN}, bson.M{})

	for _, entry := range entries {
		// Check CMDB entry is of type VLAN
		if entry.CMDBType != VLAN {
			continue
		}

		lowIP, lowFound, _ := entry.FindSysTag("LowIP")
		highIP, highFound, _ := entry.FindSysTag("HighIP")

		// Check that both of the user tags for the IPs are set
		if !lowFound && !highFound {
			continue
		}

		cidr, err := utils.IPv4RangeToCIDRRange(lowIP.Values[0], highIP.Values[0])
		system.Error("Couldn't generate CIDR for: "+entry.Label, err)

		// Remove old CMDB tags so new one can be calcualted
		_, found, index := entry.FindSysTag("CIDR")

		if found {
			entry.SysTags[index] = EntityTag{Label: "CIDR", Description: "CIDR range for this VLAN.", DataType: system.DataType_CIDR, Values: cidr}
		} else {
			entry.SysTags = append(entry.SysTags, EntityTag{Label: "CIDR", Description: "CIDR range for this VLAN.", DataType: system.DataType_CIDR, Values: cidr})
		}

		entry.UPDATE_ENTRY_Inventory()
	}
}

/*
	Given an IP, functional finds additional info and then
	creates an inventory entry for the IP.
*/
func setLocalIdentityEntry(ip string) {
	sysTags := []EntityTag{}
	usrTags := []EntityTag{}

	sysTags = append(sysTags, EntityTag{Label: "Verified", DataType: system.DataType_BOOL, Values: []string{"1"}})
	sysTags = append(sysTags, EntityTag{Label: "Identity", DataType: system.DataType_STRING, Values: []string{"local"}})
	sysTags = append(sysTags, EntityTag{Label: "IP", DataType: system.DataType_IP, Values: []string{ip}})

	for _, net := range local.GetNetworkAdapters() {
		if net.IP == ip {
			if net.MAC != "" {
				sysTags = append(sysTags, EntityTag{Label: "MAC", DataType: system.DataType_MAC, Values: []string{net.MAC}})
			}

			if net.MAC6 != "" {
				sysTags = append(sysTags, EntityTag{Label: "MAC6", DataType: system.DataType_MAC6, Values: []string{net.MAC6}})
			}

			if net.Label != "" {
				sysTags = append(sysTags, EntityTag{Label: "NetAdapter", DataType: system.DataType_STRING, Values: []string{net.Label}})
			}

			if net.IP6 != "" {
				sysTags = append(sysTags, EntityTag{Label: "IP6", DataType: system.DataType_IP6, Values: []string{net.IP6}})
			}
		}
	}

	time := []string{utils.GetDateTime().DateTime}

	serverCMDB := Entity{
		Label:       "Local-Mapper Server",
		OSILayer:    7,
		Description: "The local-mapper backend server.",
		DateSeen:    time,
		CMDBType:    SERVER,
		UsrTags:     usrTags,
		SysTags:     sysTags,
	}

	serverCMDB.InsertInventory()
}

func resolveIPConflict(action ConflictActions, ip string) {

	if action == Action_MERGE_INTO_INVENTORY {

		DELETE_ENTRY_Pending(SELECT_ENTRY_Pending(bson.M{"systags.label": "IP", "systags.values": ip}, bson.M{})[0])
	}

	if action == Action_MERGE_INTO_PENDING {

	}

	if action == Action_DELETE_INVENTORY_ENTRY {
		DELETE_ENTRY_Inventory(SELECT_ENTRY_Inventory(bson.M{"systags.label": "IP", "systags.values": ip}, bson.M{})[0])
	}

	if action == Action_DELETE_PENDING_ENTRY {
		DELETE_ENTRY_Pending(SELECT_ENTRY_Pending(bson.M{"systags.label": "IP", "systags.values": ip}, bson.M{})[0])
	}
}
