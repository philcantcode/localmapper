package proposition

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func INSERT_Proposition(proposition Proposition) {
	utils.Log("Attempting to INSERT_Proposition", false)
	proposition.ID = primitive.NewObjectID()
	insertResult, err := database.PropositionDB.InsertOne(context.Background(), proposition)

	utils.ErrorFatal("Couldn't INSERT_Proposition", err)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), true)
}

/* SELECT_Capability takes in a:
   1. filter bson.M{"startstr": "xyz"} select specifc values
   2. projection bson.M{"version": 1} to limit the fields returned
*/
func SELECT_Propositions(filter bson.M, projection bson.M) []Proposition {
	cursor, err := database.PropositionDB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	utils.ErrorFatal("Couldn't SELECT_Propositions", err)
	defer cursor.Close(context.Background())

	var results []Proposition

	for cursor.Next(context.Background()) {
		var prop Proposition

		err = cursor.Decode(&prop)
		utils.ErrorFatal("Couldn't decode SELECT_Propositions", err)

		results = append(results, prop)
	}

	return results
}

/*
Status:
0 = Open
1 = Complete
2 = Disabled */
func UPDATE_Proposition(proposition Proposition) {
	result, err := database.PropositionDB.ReplaceOne(context.Background(), bson.M{"_id": proposition.ID}, proposition)
	utils.ErrorFatal("Couldn't UPDATE_Proposition", err)

	utils.Log(fmt.Sprintf("UPDATE_Proposition ID: %s, Result: %d\n", proposition.ID, result.ModifiedCount), true)
}
