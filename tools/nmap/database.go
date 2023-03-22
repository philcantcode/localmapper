package nmap

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (nmap NmapRun) INSERT() string {
	system.Log("Attempting to INSERT_Nmap", false)

	insertResult, err := system.Results_Nmap_DB.InsertOne(context.Background(), nmap)

	system.Fatal("Couldn't Insert_Network_Nmap", err)
	system.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)

	return insertResult.InsertedID.(primitive.ObjectID).Hex()
}

func SELECT(filter bson.M, projection bson.M) []NmapRun {
	cursor, err := system.Results_Nmap_DB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	system.Fatal("Couldn't SELECT Nmap", err)
	defer cursor.Close(context.Background())

	results := []NmapRun{}

	for cursor.Next(context.Background()) {
		var nmap NmapRun

		err = cursor.Decode(&nmap)
		system.Fatal("Couldn't decode SELECT Nmap", err)

		results = append(results, nmap)
	}

	return results
}
