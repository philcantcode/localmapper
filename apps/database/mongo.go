package database

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connection URI
var uri string
var Client *mongo.Client

func InitMongo() {
	var err error

	utils.Log("Attempting to connect MongoDB to: "+uri, false)

	if utils.Configs["MONGO_PASSWORD_REQUIRED"] == "1" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s",
			utils.Configs["MONGO_USER"],
			utils.Configs["MONGO_PASSWORD"],
			utils.Configs["MONGO_IP"],
			utils.Configs["MONGO_PORT"])
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s",
			utils.Configs["MONGO_IP"],
			utils.Configs["MONGO_PORT"])
	}

	// Create a new client and connect to the server
	Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	utils.ErrorFatal("MongoDB couldn't make initial connection to "+uri, err)

	// Ping the primary
	err = Client.Ping(context.TODO(), readpref.Primary())
	utils.ErrorFatal("MongoDB couldn't ping "+uri, err)
	utils.Log("Successfully connected MongoDB to: "+uri, true)
}
