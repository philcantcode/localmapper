package database

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/adapters/definitions"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertNetworkNmap(nmapResult definitions.NmapRun) {
	utils.Log("Attempting to Insert_Network_Nmap", false)

	insertResult, err := nmapDB.InsertOne(context.Background(), nmapResult)

	utils.ErrorFatal("Couldn't Insert_Network_Nmap", err)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), true)
}

/* FilterNetworkNmap takes in a:
   1. filter bson.M{"startstr": "xyz"} select specifc values
   2. projection bson.M{"version": 1} to limit the fields returned
*/
func FilterNetworkNmap(filter bson.M, projection bson.M) []definitions.NmapRun {
	cursor, err := nmapDB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	utils.ErrorFatal("Couldn't FilterNetworkMap", err)
	defer cursor.Close(context.Background())

	var results []definitions.NmapRun

	for cursor.Next(context.Background()) {
		var nmapRun definitions.NmapRun

		err = cursor.Decode(&nmapRun)
		utils.ErrorFatal("Couldn't decode application.nmap.SelectAllNetworkNmap", err)

		results = append(results, nmapRun)
	}

	return results
}
