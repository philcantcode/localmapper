package cmdb

import (
	"fmt"

	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

/*
	EntryExistsSomewhere returns true if an entry exists either
	Inventory OR Pending
*/
func EntryExists_ByIP(entry Entry) bool {
	tag, exists, _ := FindSysTag("IP", entry)
	result := []Entry{}

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
	Updates Inventory database entries by IP.

	Finds the old entry result[0] and then updates the
	values to the new entry.
*/
func UpdateInventoryEntries_ByIP(entry Entry) bool {
	tag, exists, _ := FindSysTag("IP", entry)

	if !exists {
		return false
	}

	ipFilter := bson.M{
		"systags.label":  "IP",
		"systags.values": tag.Values[len(tag.Values)-1],
	}

	results := SELECT_ENTRY_Inventory(ipFilter, bson.M{})

	if len(results) == 0 {
		utils.Log(fmt.Sprintf("No match for (inventory): %s\n", tag.Values[len(tag.Values)-1]), false)
		return false
	}

	if len(results) > 1 { // Too many results returned, database corrupt
		utils.ErrorForceFatal(
			fmt.Sprintf(
				"While executing UpdateInventoryEntries the number of matched results > 1\n\nEntry: %+v\n\nMatched Cases: %+v", entry, results))
	}

	utils.Log(fmt.Sprintf("Match (Inventory): len: %d, IP: %+v\n", len(results), results), false)

	for _, newTag := range entry.SysTags {
		_, found, i := FindSysTag(newTag.Label, results[0])

		if found {
			results[0].SysTags[i].Values = newTag.Values
		} else {
			results[0].SysTags = append(results[0].SysTags, newTag)
		}
	}

	for _, newTag := range entry.UsrTags {
		_, found, i := FindUsrTag(newTag.Label, results[0])

		if found {
			results[0].UsrTags[i].Values = newTag.Values
		} else {
			results[0].UsrTags = append(results[0].UsrTags, newTag)
		}
	}

	results[0].DateSeen = append(results[0].DateSeen, entry.DateSeen...)

	utils.Log(fmt.Sprintf("Compartive update made for (inventory): %v\n", results[0].ID), false)
	UPDATE_ENTRY_Inventory(results[0])

	return true
}

/*
	Updates Pending database entries

	Finds the old entry result[0] and then updates the
	values to the new entry.
*/
func UpdatePendingEntries_ByIP(entry Entry) bool {
	tag, exists, _ := FindSysTag("IP", entry)

	if !exists {
		return false
	}

	ipFilter := bson.M{
		"systags.label":  "IP",
		"systags.values": tag.Values[len(tag.Values)-1],
	}

	results := SELECT_ENTRY_Pending(ipFilter, bson.M{})

	if len(results) == 0 {
		utils.Log(fmt.Sprintf("No match for (pending): %s\n", tag.Values[len(tag.Values)-1]), false)
		return false
	}

	// Too many results returned, database corrupt
	if len(results) > 1 {
		utils.ErrorForceFatal(
			fmt.Sprintf(
				"While executing UpdateInventoryEntries the number of matched results > 1\n\nEntry: %+v\n\nMatched Cases: %+v", entry, results))
	}

	utils.Log(fmt.Sprintf("Match (Pending): len: %d, IP: %+v\n", len(results), results), false)

	for _, newTag := range entry.SysTags {
		_, found, i := FindSysTag(newTag.Label, results[0])

		if found {
			results[0].SysTags[i].Values = newTag.Values
		} else {
			results[0].SysTags = append(results[0].SysTags, newTag)
		}
	}

	for _, newTag := range entry.UsrTags {
		_, found, i := FindUsrTag(newTag.Label, results[0])

		if found {
			results[0].UsrTags[i].Values = newTag.Values
		} else {
			results[0].UsrTags = append(results[0].UsrTags, newTag)
		}
	}

	results[0].Label = entry.Label
	results[0].Desc = entry.Desc
	results[0].CMDBType = entry.CMDBType
	results[0].OSILayer = entry.OSILayer
	results[0].DateSeen = append(results[0].DateSeen, entry.DateSeen...)

	utils.Log(fmt.Sprintf("Compartive update made for (pending): %v\n", results[0].ID), false)
	UPDATE_ENTRY_Pending(results[0])

	return true
}
