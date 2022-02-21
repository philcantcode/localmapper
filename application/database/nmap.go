package database

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/adapters/blueprint"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertNetworkNmap(nmapResult blueprint.NmapRun) {
	utils.Log("Attempting to Insert_Network_Nmap", false)

	nmapCollection := client.Database("Network").Collection("Nmap")
	insertResult, err := nmapCollection.InsertOne(context.Background(), nmapResult)

	utils.ErrorFatal("Couldn't Insert_Network_Nmap", err)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), true)
}

func FilterNetworkMap(filter bson.M) {
	nmapCollection := client.Database("Network").Collection("Nmap")
	cursor, err := nmapCollection.Find(context.Background(), filter)
	utils.ErrorFatal("Couldn't SelectAllNetworkNmap", err)

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var nmapRun blueprint.NmapRun

		err = cursor.Decode(&nmapRun)
		utils.ErrorFatal("Couldn't decode application.nmap.SelectAllNetworkNmap", err)

		fmt.Println(nmapRun.StartStr)
	}
}
