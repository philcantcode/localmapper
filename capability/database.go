package capability

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func INSERT_Capability(capability Capability) {
	utils.Log("Attempting to INSERT_Capability", false)

	capability.ID = primitive.NewObjectID()
	insertResult, err := database.Core_Capability_DB.InsertOne(context.Background(), capability)

	utils.ErrorFatal("Couldn't INSERT_Capability", err)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)
}

func UPDATE_Capability(capability Capability) {
	utils.Log(fmt.Sprintf("Attempting to UPDATE_Capability (ID: %d)\n", capability.ID), true)

	result, err := database.Core_Capability_DB.ReplaceOne(context.Background(), bson.M{"id": capability.ID}, capability)

	utils.ErrorFatal("Couldn't UPDATE_Capability", err)
	utils.Log(fmt.Sprintf("New Update made: %b", result.ModifiedCount), false)
}

/* SELECT_Capability takes in a:
   1. filter bson.M{"startstr": "xyz"} select specifc values
   2. projection bson.M{"version": 1} to limit the fields returned
*/
func SELECT_Capability(filter bson.M, projection bson.M) []Capability {
	cursor, err := database.Core_Capability_DB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	utils.ErrorFatal("Couldn't SELECT_Capability", err)
	defer cursor.Close(context.Background())

	var results []Capability

	for cursor.Next(context.Background()) {
		var cap Capability

		err = cursor.Decode(&cap)
		utils.ErrorFatal("Couldn't decode SELECT_Capability", err)

		results = append(results, cap)
	}

	return results
}
