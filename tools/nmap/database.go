package nmap

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (nmap NmapRun) Insert() string {
	system.Log("Attempting to INSERT_Nmap", false)

	insertResult, err := system.Results_Nmap_DB.InsertOne(context.Background(), nmap)

	system.Fatal("Couldn't Insert_Network_Nmap", err)
	system.Log(fmt.Sprintf("New Insert at: %s", insertResult), true)

	return insertResult.InsertedID.(primitive.ObjectID).Hex()
}
