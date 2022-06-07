package cmdb

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func INSERT_ENTRY_Inventory(cmdb Entry) {
	utils.Log("Attempting to INSERT_CMDB", false)

	cmdb.ID = primitive.NewObjectID()
	insertResult, err := database.CMDB_Inventory_DB.InsertOne(context.Background(), cmdb)

	utils.ErrorFatal("Couldn't INSERT_CMDB", err)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)
}

func INSERT_ENTRY_Pending(entry Entry) {
	utils.Log("Attempting to INSERT_ENTRY_Pending", false)
	matchedEntry, entryExists := matchedEntryExists(entry, true)

	if entryExists {
		utils.Log("Entry in INSERT_ENTRY_Pending is already registered", false)

		updateEntry(matchedEntry, entry, true)
		return
	}

	entry.ID = primitive.NewObjectID()
	insertResult, err := database.CMDB_Pending_DB.InsertOne(context.Background(), entry)

	utils.ErrorFatal("Couldn't INSERT_ENTRY_Pending", err)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)
}

func SELECT_ENTRY_Inventory(filter bson.M, projection bson.M) []Entry {
	cursor, err := database.CMDB_Inventory_DB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	utils.ErrorFatal("Couldn't SELECT_CMDBItem", err)
	defer cursor.Close(context.Background())

	results := []Entry{}

	for cursor.Next(context.Background()) {
		var cmdb Entry

		err = cursor.Decode(&cmdb)
		utils.ErrorFatal("Couldn't decode SELECT_CMDBItem", err)

		results = append(results, cmdb)
	}

	return results
}

func SELECT_ENTRY_Pending(filter bson.M, projection bson.M) []Entry {
	cursor, err := database.CMDB_Pending_DB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	utils.ErrorFatal("Couldn't SELECT_ENTRY_Pending", err)
	defer cursor.Close(context.Background())

	results := []Entry{}

	for cursor.Next(context.Background()) {
		var cmdb Entry

		err = cursor.Decode(&cmdb)
		utils.ErrorFatal("Couldn't decode SELECT_ENTRY_Pending", err)

		results = append(results, cmdb)
	}

	return results
}

func UPDATE_ENTRY_Inventory(cmdb Entry) {
	result, err := database.CMDB_Inventory_DB.ReplaceOne(context.Background(), bson.M{"_id": cmdb.ID}, cmdb)
	utils.ErrorFatal("Couldn't UPDATE_ENTRY_Inventory", err)

	utils.Log(fmt.Sprintf("UPDATE_ENTRY_Inventory ID: %s, Result: %d\n", cmdb.ID, result.ModifiedCount), false)
}

func UPDATE_ENTRY_Pending(cmdb Entry) {
	result, err := database.CMDB_Pending_DB.ReplaceOne(context.Background(), bson.M{"_id": cmdb.ID}, cmdb)
	utils.ErrorFatal("Couldn't UPDATE_ENTRY_Pending", err)

	utils.Log(fmt.Sprintf("UPDATE_ENTRY_Pending ID: %s, Result: %d\n", cmdb.ID, result.ModifiedCount), false)
}

func DELETE_ENTRY_Pending(entry Entry) {
	utils.Log("Attempting to DELETE_ENTRY_Pending", false)

	insertResult, err := database.CMDB_Pending_DB.DeleteOne(context.Background(), bson.M{"_id": entry.ID})

	utils.ErrorFatal("Couldn't DELETE_ENTRY_Pending", err)
	utils.Log(fmt.Sprintf("New Delete count: %d", insertResult.DeletedCount), false)
}

/*
matchedEntryExists checks various values to see if they already exist
in the pending database

Match against pendingDB = {true = PENDING}{false = INVENTORY}
*/
func matchedEntryExists(entry Entry, pendingDB bool) (Entry, bool) {
	filter := bson.M{}

	for _, tag := range entry.SysTags {
		switch tag.DataType {
		case utils.IP:
			filter["systags.label"] = "IP"
			filter["systags.values"] = tag.Values[0]

		case utils.IP6:
			filter["systags.label"] = "IP6"
			filter["systags.values"] = tag.Values[0]

		case utils.MAC:
			filter["systags.label"] = "MAC"
			filter["systags.values"] = tag.Values[0]

		case utils.MAC6:
			filter["systags.label"] = "MAC6"
			filter["systags.values"] = tag.Values[0]
		}

		switch tag.Label {
		case "HostName":
			filter["systags.label"] = "HostName"
			filter["systags.values"] = tag.Values[0]
		}
	}

	result := []Entry{}

	// Select which database
	if pendingDB {
		result = SELECT_ENTRY_Pending(filter, bson.M{})
	} else {
		result = SELECT_ENTRY_Inventory(filter, bson.M{})
	}

	//TODO: Case where the new entry matched multiple entries in DB
	if len(result) > 1 {
		utils.ErrorContextLog("Multiple matched entries have been returned, handle with proposition", true)
	}

	return result[0], len(result) != 0
}

/* Updates old entry to new entry variables & merges variables where appropriate*/
func updateEntry(oldEntry Entry, newEntry Entry, pendingDB bool) {
	// For System Tags
	for _, tag := range newEntry.SysTags {
		_, found, index := FindSysTag(tag.Label, oldEntry)

		// Tag [Label] already exists, overwrite values
		// TODO: In future we should loop over each value and append where necessary (last most up to date)
		if found {
			oldEntry.SysTags[index].Values = tag.Values
		} else {
			oldEntry.SysTags = append(oldEntry.SysTags, tag)
		}
	}

	// For User Tags
	for _, tag := range newEntry.UsrTags {
		_, found, index := FindUsrTag(tag.Label, oldEntry)

		// Tag [Label] already exists, overwrite values
		// TODO: In future we should loop over each value and append where necessary (last most up to date)
		if found {
			oldEntry.UsrTags[index].Values = tag.Values
		} else {
			oldEntry.UsrTags = append(oldEntry.UsrTags, tag)
		}
	}

	// For metadata
	oldEntry.Label = newEntry.Label
	oldEntry.Desc = newEntry.Desc
	oldEntry.CMDBType = newEntry.CMDBType
	oldEntry.DateSeen = append(oldEntry.DateSeen, newEntry.DateSeen...)
	oldEntry.OSILayer = newEntry.OSILayer

	utils.Log(fmt.Sprintf("Compartive update made for: %v\n", oldEntry.ID), false)
	UPDATE_ENTRY_Pending(oldEntry)
}
