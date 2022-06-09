package cookbook

import (
	"fmt"

	"github.com/philcantcode/localmapper/capability"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
	ExecuteCookbook runs a passed cookbook
*/
func ExecuteCookbook(book Cookbook, entryID primitive.ObjectID) {
	utils.Log("Attempting to execute cookbook: "+book.Label, false)

	capList := capability.SELECT_Capability(bson.M{}, bson.M{})

	for _, cci := range book.CCIs {
		// Select the entry and capability given an id and cci
		entries := cmdb.SELECT_ENTRY_Joined(bson.M{"_id": entryID}, bson.M{})
		caps := capability.SELECT_Capability(bson.M{"cci": cci}, bson.M{})

		// Ensure only 1 capability returned
		if len(caps) != 1 {
			utils.ErrorContextLog(
				fmt.Sprintf(
					"Incorrect number (%d) of returned for CCI: %s",
					len(caps), cci), true)
			return
		}

		// Ensure only 1 entry returned
		if len(entries) != 1 {
			utils.ErrorContextLog(
				fmt.Sprintf(
					"Incorrect number (%d) of returned for entries: %s",
					len(entries), entryID), true)
			return
		}

		isMatch, cap := capability.MatchEntryToCapability(caps[0], entries[0])

		if isMatch {
			utils.Log("Matched entry found, performing capability", false)
			capability.ExecuteCapability(cap)
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
						utils.ErrorContextLog(
							fmt.Sprintf(
								"Incorrect number (%d) of returned for entries: %s",
								len(entries), entryID), true)
						return
					}

					isMatch, cap := capability.MatchEntryToCapability(cap, entries[0])

					if isMatch {
						utils.Log("Matched entry found (using searchKey), performing capability", false)
						capability.ExecuteCapability(cap)
					}
				}
			}
		}
	}
}
