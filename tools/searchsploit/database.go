package searchsploit

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (result ExploitDB) Insert() string {
	system.Log("Attempting to INSERT_Searchsploit", false)
	insertResult, err := system.Results_Searchsploit_DB.InsertOne(context.Background(), result)

	system.Fatal("Couldn't INSERT_Searchsploit", err)
	system.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)

	return insertResult.InsertedID.(primitive.ObjectID).Hex()
}

func Select(filter bson.M, projection bson.M) []ExploitDB {
	cursor, err := system.Results_Searchsploit_DB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	system.Fatal("Couldn't SELECT_Searchsploit", err)
	defer cursor.Close(context.Background())

	var results []ExploitDB

	for cursor.Next(context.Background()) {
		var exploitDB ExploitDB

		err = cursor.Decode(&exploitDB)
		system.Fatal("Couldn't decode SELECT_Searchsploit", err)

		results = append(results, exploitDB)
	}

	return results
}

func Delete(filter bson.M) {
	system.Log("Attempting to delete_searchsploit", false)

	insertResult, err := system.Results_Searchsploit_DB.DeleteMany(context.Background(), filter)

	system.Fatal("Couldn't delete_searchsploit", err)
	system.Log(fmt.Sprintf("New Delete count: %d", insertResult.DeletedCount), false)
}
