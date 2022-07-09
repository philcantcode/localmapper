package cmdb

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (entry Entity) InsertInventory() {
	system.Log("Attempting to INSERT_ENTRY_Inventory", false)

	entry.ID = primitive.NewObjectID()
	insertResult, err := system.CMDB_Inventory_DB.InsertOne(context.Background(), entry)

	system.Fatal("Couldn't INSERT_ENTRY_Inventory", err)
	system.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)
}

func (entry Entity) InsertPending() {
	system.Log("Attempting to insert_ENTRY_Pending", false)

	// Otherwise, add it to pending
	entry.ID = primitive.NewObjectID()
	insertResult, err := system.CMDB_Pending_DB.InsertOne(context.Background(), entry)

	system.Fatal("Couldn't insert_ENTRY_Pending", err)
	system.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)
}

/*
	SELECT_ENTRY_Inventory returns an array of Entry.

	Array len() = 0 if none match
*/
func SELECT_ENTRY_Inventory(filter bson.M, projection bson.M) []Entity {
	cursor, err := system.CMDB_Inventory_DB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	system.Fatal("Couldn't SELECT_CMDBItem", err)
	defer cursor.Close(context.Background())

	results := []Entity{}

	for cursor.Next(context.Background()) {
		var cmdb Entity

		err = cursor.Decode(&cmdb)
		system.Fatal("Couldn't decode SELECT_CMDBItem", err)

		results = append(results, cmdb)
	}

	return results
}

func SELECT_ENTRY_Pending(filter bson.M, projection bson.M) []Entity {
	cursor, err := system.CMDB_Pending_DB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	system.Fatal("Couldn't SELECT_ENTRY_Pending", err)
	defer cursor.Close(context.Background())

	results := []Entity{}

	for cursor.Next(context.Background()) {
		var cmdb Entity

		err = cursor.Decode(&cmdb)
		system.Fatal("Couldn't decode SELECT_ENTRY_Pending", err)

		results = append(results, cmdb)
	}

	return results
}

func SELECT_Entities_Joined(filter bson.M, projection bson.M) []Entity {
	results := []Entity{}

	results = append(results, SELECT_ENTRY_Inventory(filter, projection)...)
	results = append(results, SELECT_ENTRY_Pending(filter, projection)...)

	return results
}

func (cmdb Entity) UPDATE_ENTRY_Inventory() {
	result, err := system.CMDB_Inventory_DB.ReplaceOne(context.Background(), bson.M{"_id": cmdb.ID}, cmdb)
	system.Fatal("Couldn't UPDATE_ENTRY_Inventory", err)

	system.Log(fmt.Sprintf("UPDATE_ENTRY_Inventory ID: %s, Result: %d", cmdb.ID, result.ModifiedCount), false)
}

func (cmdb Entity) UPDATE_ENTRY_Pending() {
	result, err := system.CMDB_Pending_DB.ReplaceOne(context.Background(), bson.M{"_id": cmdb.ID}, cmdb)
	system.Fatal("Couldn't UPDATE_ENTRY_Pending", err)

	system.Log(fmt.Sprintf("UPDATE_ENTRY_Pending ID: %s, Result: %d", cmdb.ID, result.ModifiedCount), false)
}

func DELETE_ENTRY_Pending(entry Entity) {
	system.Log("Attempting to DELETE_ENTRY_Pending", false)

	insertResult, err := system.CMDB_Pending_DB.DeleteOne(context.Background(), bson.M{"_id": entry.ID})

	system.Fatal("Couldn't DELETE_ENTRY_Pending", err)
	system.Log(fmt.Sprintf("New Delete count: %d", insertResult.DeletedCount), false)
}

func DELETE_ENTRY_Inventory(entry Entity) {
	system.Log("Attempting to DELETE_ENTRY_Inventory", false)

	insertResult, err := system.CMDB_Inventory_DB.DeleteOne(context.Background(), bson.M{"_id": entry.ID})

	system.Fatal("Couldn't DELETE_ENTRY_Inventory", err)
	system.Log(fmt.Sprintf("New Delete count: %d", insertResult.DeletedCount), false)
}
