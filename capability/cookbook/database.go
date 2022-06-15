package cookbook

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SELECT_Cookbook(filter bson.M, projection bson.M) []Cookbook {
	cursor, err := system.Core_Cookbooks_DB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	system.Fatal("Couldn't SELECT_Cookbook", err)
	defer cursor.Close(context.Background())

	var results []Cookbook

	for cursor.Next(context.Background()) {
		var book Cookbook

		err = cursor.Decode(&book)
		system.Fatal("Couldn't decode SELECT_Cookbook", err)

		results = append(results, book)
	}

	return results
}

func INSERT_Cookbook(book Cookbook) {
	system.Log("Attempting to INSERT_Cookbook", false)

	book.ID = primitive.NewObjectID()
	insertResult, err := system.Core_Cookbooks_DB.InsertOne(context.Background(), book)

	system.Fatal("Couldn't INSERT_Cookbook", err)
	system.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)
}

func UPDATE(book Cookbook) {
	result, err := system.Core_Cookbooks_DB.ReplaceOne(context.Background(), bson.M{"_id": book.ID}, book)
	system.Fatal("Couldn't UPDATE Cookbook", err)

	system.Log(fmt.Sprintf("UPDATE Cookbook ID: %s, Result: %d", book.ID, result.ModifiedCount), false)
}
