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

func INSERT_ENTRY_Inventory(entry Entry) {
	utils.Log("Attempting to INSERT_ENTRY_Inventory", false)

	entry.ID = primitive.NewObjectID()
	insertResult, err := database.CMDB_Inventory_DB.InsertOne(context.Background(), entry)

	utils.ErrorFatal("Couldn't INSERT_ENTRY_Inventory", err)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)
}

func INSERT_ENTRY_Pending(entry Entry) {
	utils.Log("Attempting to INSERT_ENTRY_Pending", false)

	// Otherwise, add it to pending
	entry.ID = primitive.NewObjectID()
	insertResult, err := database.CMDB_Pending_DB.InsertOne(context.Background(), entry)

	utils.ErrorFatal("Couldn't INSERT_ENTRY_Pending", err)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)
}

/*
	SELECT_ENTRY_Inventory returns an array of Entry.

	Array len() = 0 if none match
*/
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
