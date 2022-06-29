package cmdb

import (
	"fmt"
	"strconv"
	"time"

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

	results := SELECT_ENTRY_Joined(ipFilter, bson.M{})

	if len(results) == 0 {
		system.Log(fmt.Sprintf("No match for (inventory): %s", tag.Values[len(tag.Values)-1]), false)
		return false
	}

	if len(results) > 1 { // Too many results returned, database corrupt
		system.Warning(
			fmt.Sprintf(
				"While executing UpdateInventoryEntries the number of matched results > 1\n\nEntry: %+v\n\nMatched Cases: %+v", entry, results), true)
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
		_, found, i := FindUsrTag(newTag.Label, results[0])

		if found {
			results[0].UsrTags[i].Values = joinTagGroups(newTag.Label, results[0].UsrTags[i].Values, newTag.Values)
		} else {
			results[0].UsrTags = append(results[0].UsrTags, newTag)
		}
	}

	results[0].DateSeen = append(results[0].DateSeen, entry.DateSeen...)

	system.Log(fmt.Sprintf("Compartive update made: %v", results[0].ID), false)
	UPDATE_ENTRY_Inventory(results[0])

	// Only update the metadata for the pending entry
	results[0].Label = entry.Label
	results[0].Description = entry.Description
	results[0].CMDBType = entry.CMDBType
	results[0].OSILayer = entry.OSILayer

	UPDATE_ENTRY_Pending(results[0])

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
	if len(SELECT_ENTRY_Inventory(bson.M{}, bson.M{})) > 0 {
		return
	}

	vlan1 := SELECT_ENTRY_Inventory(bson.M{"label": "Private Range 1", "desc": "Default VLAN"}, bson.M{})

	if len(vlan1) == 0 {
		lowIP := EntityTag{Label: "LowIP", DataType: system.DataType_IP_RANGE_LOW, Values: []string{"10.0.0.0"}}
		highIP := EntityTag{Label: "HighIP", DataType: system.DataType_IP_RANGE_HIGH, Values: []string{"10.255.255.255"}}
		sysDefault := EntityTag{Label: "SysDefault", DataType: system.DataType_BOOL, Values: []string{"1"}}

		newVlan := Entity{Label: "Private Range 1", Description: "Default VLAN", CMDBType: VLAN, OSILayer: 2, DateSeen: []string{utils.Now()}, SysTags: []EntityTag{lowIP, highIP, sysDefault}}
		insert_ENTRY_Inventory(newVlan)
	}

	vlan2 := SELECT_ENTRY_Inventory(bson.M{"label": "Private Range 2", "desc": "Default VLAN"}, bson.M{})

	if len(vlan2) == 0 {
		lowIP := EntityTag{Label: "LowIP", DataType: system.DataType_IP_RANGE_LOW, Values: []string{"172.16.0.0"}}
		highIP := EntityTag{Label: "HighIP", DataType: system.DataType_IP_RANGE_HIGH, Values: []string{"172.31.255.255"}}
		sysDefault := EntityTag{Label: "SysDefault", DataType: system.DataType_BOOL, Values: []string{"1"}}

		newVlan := Entity{Label: "Private Range 2", Description: "Default VLAN", CMDBType: VLAN, OSILayer: 2, DateSeen: []string{utils.Now()}, SysTags: []EntityTag{lowIP, highIP, sysDefault}}
		insert_ENTRY_Inventory(newVlan)
	}

	vlan3 := SELECT_ENTRY_Inventory(bson.M{"label": "Private Range 3", "desc": "Default VLAN"}, bson.M{})

	if len(vlan3) == 0 {
		lowIP := EntityTag{Label: "LowIP", DataType: system.DataType_IP_RANGE_LOW, Values: []string{"192.168.0.0"}}
		highIP := EntityTag{Label: "HighIP", DataType: system.DataType_IP_RANGE_HIGH, Values: []string{"192.168.255.255"}}
		sysDefault := EntityTag{Label: "SysDefault", DataType: system.DataType_BOOL, Values: []string{"1"}}

		newVlan := Entity{Label: "Private Range 3", Description: "Default VLAN", CMDBType: VLAN, OSILayer: 2, DateSeen: []string{utils.Now()}, SysTags: []EntityTag{lowIP, highIP, sysDefault}}
		insert_ENTRY_Inventory(newVlan)
	}

	vlan4 := SELECT_ENTRY_Inventory(bson.M{"label": "Test Home", "desc": "Test VLAN"}, bson.M{})

	if len(vlan4) == 0 {
		lowIP := EntityTag{Label: "LowIP", DataType: system.DataType_IP_RANGE_LOW, Values: []string{"192.168.1.0"}}
		highIP := EntityTag{Label: "HighIP", DataType: system.DataType_IP_RANGE_HIGH, Values: []string{"192.168.1.255"}}
		sysDefault := EntityTag{Label: "SysDefault", DataType: system.DataType_BOOL, Values: []string{"1"}}

		newVlan := Entity{Label: "Test Home", Description: "Test VLAN", CMDBType: VLAN, OSILayer: 2, DateSeen: []string{utils.Now()}, SysTags: []EntityTag{lowIP, highIP, sysDefault}}
		insert_ENTRY_Inventory(newVlan)
	}

	vlan5 := SELECT_ENTRY_Inventory(bson.M{"label": "Olivers Home", "desc": "Test VLAN"}, bson.M{})

	if len(vlan5) == 0 {
		lowIP := EntityTag{Label: "LowIP", DataType: system.DataType_IP_RANGE_LOW, Values: []string{"192.168.0.0"}}
		highIP := EntityTag{Label: "HighIP", DataType: system.DataType_IP_RANGE_HIGH, Values: []string{"192.168.0.255"}}
		sysDefault := EntityTag{Label: "SysDefault", DataType: system.DataType_BOOL, Values: []string{"1"}}

		newVlan := Entity{Label: "Olivers Home", Description: "Test VLAN", CMDBType: VLAN, OSILayer: 2, DateSeen: []string{utils.Now()}, SysTags: []EntityTag{lowIP, highIP, sysDefault}}
		insert_ENTRY_Inventory(newVlan)
	}
}

/*
	UpdateOrInsert either updates the entity by IP if found
	or inserts a new entity
*/
func (entry Entity) UpdateOrInsert() {
	// Insert to pending or update both DBs
	if EntryExists_ByIP(entry) {
		entryUpdateSuccess := updateEntriesTags_ByIP(entry)

		if !entryUpdateSuccess {
			system.Warning("Couldn't update inventory or pending in CMDB", true)
		}
	} else {
		insert_ENTRY_Pending(entry)
	}
}
