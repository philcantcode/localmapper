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

func INSERT_ENTRY(cmdb Entry) {
	utils.Log("Attempting to INSERT_CMDB", false)

	cmdb.ID = primitive.NewObjectID()
	insertResult, err := database.CMDB_Inventory_DB.InsertOne(context.Background(), cmdb)

	utils.ErrorFatal("Couldn't INSERT_CMDB", err)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)
}

func SELECT_ENTRY(filter bson.M, projection bson.M) []Entry {
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

func UPDATE_ENTRY(cmdb Entry) {
	result, err := database.CMDB_Inventory_DB.ReplaceOne(context.Background(), bson.M{"_id": cmdb.ID}, cmdb)
	utils.ErrorFatal("Couldn't UPDATE_CMDB", err)

	fmt.Printf("%+v\n", cmdb)
	utils.Log(fmt.Sprintf("UPDATE_CMDB ID: %s, Result: %d\n", cmdb.ID, result.ModifiedCount), false)
}
