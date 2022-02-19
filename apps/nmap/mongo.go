package nmap

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/apps/database"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

var nmapCollection *mongo.Collection

func init() {

}

func MongoInsert(v interface{}) {
	nmapCollection = database.Client.Database("Network").Collection("Nmap")
	insertResult, err := nmapCollection.InsertOne(context.TODO(), v)

	utils.ErrorLog("Couldn't Insert (mongodb)", err, true)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), true)
}
