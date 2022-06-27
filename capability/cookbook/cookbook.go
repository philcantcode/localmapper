package cookbook

import (
	"fmt"

	"github.com/philcantcode/localmapper/capability"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
	ExecuteOnEntry runs a cookbook against a given entry
*/
func (book Cookbook) ExecuteOnEntry(entryID primitive.ObjectID) {
	system.Log("Attempting to execute cookbook: "+book.Label, true)
	capsInBook := []string{} // Keep track of capabilities already run so don't run them twice in 1 cookbook

	capList := capability.SELECT_Capability(bson.M{}, bson.M{})

	for _, cci := range book.CCIs {
		// Select the entry and capability given an id and cci
		entries := cmdb.SELECT_ENTRY_Joined(bson.M{"_id": entryID}, bson.M{})
		caps := capability.SELECT_Capability(bson.M{"cci": cci}, bson.M{})

		// Ensure only 1 capability returned
		if len(caps) != 1 {
			system.Warning(
				fmt.Sprintf(
					"Incorrect number (%d) of returned for CCI: %s",
					len(caps), cci), true)
			return
		}

		// Ensure only 1 entry returned
		if len(entries) != 1 {
			system.Warning(
				fmt.Sprintf(
					"Incorrect number (%d) of returned for entries: %s",
					len(entries), entryID), true)
			return
		}

		isMatch, cap := caps[0].CheckCompatability(entries[0])

		if isMatch && !utils.ArrayContains(cap.ID.Hex(), capsInBook) {
			system.Log(fmt.Sprintf("Executing capability [%s] against [%s]", caps[0].Label, entries[0].Label), true)
			cap.QueueCapability()
			capsInBook = append(capsInBook, cap.ID.Hex())
		} else {
			system.Log(fmt.Sprintf("Can't execute capability [%s] against [%s], not a match", caps[0].Label, entries[0].Label), true)
		}
	}

	/*
		For each searchKey, loop the capabilities and find one
		that can produce it as a result, e.g., stealth scan can
		produce Ports.
	*/
	for _, key := range book.SearchKeys {
		for _, cap := range capList {
			// Filter out CCIs which have already been run in previous step
			if !utils.ArrayContains(cap.CCI, book.CCIs) {
				if utils.ArrayContains(key, cap.ResultTags) {
					// Reselect, last pass through might of updated tags
					entries := cmdb.SELECT_ENTRY_Joined(bson.M{"_id": entryID}, bson.M{})

					// Ensure only 1 entry returned
					if len(entries) != 1 {
						system.Warning(
							fmt.Sprintf(
								"Incorrect number (%d) of returned for entries: %s",
								len(entries), entryID), true)
						return
					}

					isMatch, cap := cap.CheckCompatability(entries[0])

					if isMatch && !utils.ArrayContains(cap.ID.Hex(), capsInBook) {
						system.Log(fmt.Sprintf("Executing capability [%s] against [%s]", cap.Label, entries[0].Label), true)
						cap.QueueCapability()
						capsInBook = append(capsInBook, cap.ID.Hex())
					}
				}
			}
		}
	}
}
