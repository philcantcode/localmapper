package database

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/adapters/blueprint"
	"github.com/philcantcode/localmapper/utils"
)

func InsertNetworkNmap(nmapResult blueprint.NmapRun) {
	utils.Log("Attempting to Insert_Network_Nmap", false)

	nmapCollection := client.Database("Network").Collection("Nmap")
	insertResult, err := nmapCollection.InsertOne(context.TODO(), nmapResult)

	utils.ErrorFatal("Couldn't Insert_Network_Nmap", err)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), true)
}
