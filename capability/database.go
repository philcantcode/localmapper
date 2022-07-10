package capability

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func INSERT_Capability(capability Capability) {
	system.Log("Attempting to INSERT_Capability", false)

	capability.ID = primitive.NewObjectID()
	insertResult, err := system.Core_Capability_DB.InsertOne(context.Background(), capability)

	system.Fatal("Couldn't INSERT_Capability", err)
	system.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)
}

func (capability *Capability) UPDATE_Capability() {
	system.Log(fmt.Sprintf("Attempting to UPDATE_Capability (ID: %d)", capability.ID), true)

	result, err := system.Core_Capability_DB.ReplaceOne(context.Background(), bson.M{"id": capability.ID}, capability)

	system.Fatal("Couldn't UPDATE_Capability", err)
	system.Log(fmt.Sprintf("New Update made: %b", result.ModifiedCount), false)
}

/* SELECT_Capability takes in a:
   1. filter bson.M{"startstr": "xyz"} select specifc values
   2. projection bson.M{"version": 1} to limit the fields returned
*/
func SELECT_Capability(filter bson.M, projection bson.M) []Capability {
	cursor, err := system.Core_Capability_DB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	system.Fatal("Couldn't SELECT_Capability", err)
	defer cursor.Close(context.Background())

	var results []Capability

	for cursor.Next(context.Background()) {
		var cap Capability

		err = cursor.Decode(&cap)
		system.Fatal("Couldn't decode SELECT_Capability: ", err)

		results = append(results, cap)
	}

	return results
}

func DELETE_Capability(filter bson.M) {
	system.Log("Attempting to DELETE_Capability", false)

	insertResult, err := system.Core_Capability_DB.DeleteOne(context.Background(), filter)

	system.Fatal("Couldn't DELETE_Capability", err)
	system.Log(fmt.Sprintf("New Delete count: %d", insertResult.DeletedCount), false)
}

func Restore() {
	system.Log("Restoring capabilities to factory defaults", true)

	DELETE_Capability(bson.M{})
	system.Core_Capability_DB.Drop(context.Background()) // Drop capabilities

	Init() // Restore capabilities
}
