package cookbook

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SELECT_Cookbook(filter bson.M, projection bson.M) []Cookbook {
	cursor, err := database.Core_Cookbooks_DB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	utils.ErrorFatal("Couldn't SELECT_Cookbook", err)
	defer cursor.Close(context.Background())

	var results []Cookbook

	for cursor.Next(context.Background()) {
		var book Cookbook

		err = cursor.Decode(&book)
		utils.ErrorFatal("Couldn't decode SELECT_Cookbook", err)

		results = append(results, book)
	}

	return results
}

func INSERT_Cookbook(book Cookbook) {
	utils.Log("Attempting to INSERT_Cookbook", false)

	book.ID = primitive.NewObjectID()
	insertResult, err := database.Core_Cookbooks_DB.InsertOne(context.Background(), book)

	utils.ErrorFatal("Couldn't INSERT_Cookbook", err)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)
}
