package nmap

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func INSERT_Nmap(nmapResult NmapRun) {
	system.Log("Attempting to INSERT_Nmap", false)

	insertResult, err := system.Results_Nmap_DB.InsertOne(context.Background(), nmapResult)

	system.Fatal("Couldn't Insert_Network_Nmap", err)
	system.Log(fmt.Sprintf("New Insert at: %s", insertResult), true)
}

/* SELECT_Nmap takes in a:
   1. filter bson.M{"startstr": "xyz"} select specifc values
   2. projection bson.M{"version": 1} to limit the fields returned
*/
func SELECT_Nmap(filter bson.M, projection bson.M) []NmapRun {
	cursor, err := system.Results_Nmap_DB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	system.Fatal("Couldn't FilterNetworkMap", err)
	defer cursor.Close(context.Background())

	var results []NmapRun

	for cursor.Next(context.Background()) {
		var nmapRun NmapRun

		err = cursor.Decode(&nmapRun)
		system.Fatal("Couldn't decode application.nmap.SelectAllNetworkNmap", err)

		results = append(results, nmapRun)
	}

	return results
}
