package proposition

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (proposition Proposition) Insert() {
	system.Log("Attempting to INSERT_Proposition", false)
	proposition.ID = primitive.NewObjectID()
	insertResult, err := system.Core_Proposition_DB.InsertOne(context.Background(), proposition)

	system.Fatal("Couldn't INSERT_Proposition", err)
	system.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)
}

/* SELECT_Capability takes in a:
   1. filter bson.M{"startstr": "xyz"} select specifc values
   2. projection bson.M{"version": 1} to limit the fields returned
*/
func SELECT_Propositions(filter bson.M, projection bson.M) []Proposition {
	cursor, err := system.Core_Proposition_DB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	system.Fatal("Couldn't SELECT_Propositions", err)
	defer cursor.Close(context.Background())

	results := []Proposition{}

	for cursor.Next(context.Background()) {
		var prop Proposition

		err = cursor.Decode(&prop)
		system.Fatal("Couldn't decode SELECT_Propositions", err)

		results = append(results, prop)
	}

	return results
}

/*
Status:
0 = Open
1 = Complete
2 = Disabled */
func (proposition Proposition) Update() {
	result, err := system.Core_Proposition_DB.ReplaceOne(context.Background(), bson.M{"_id": proposition.ID}, proposition)
	system.Fatal("Couldn't UPDATE_Proposition", err)

	system.Log(fmt.Sprintf("UPDATE_Proposition ID: %s, Result: %d", proposition.ID, result.ModifiedCount), false)
}
