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
var client *mongo.Client
var Results_Nmap_DB *mongo.Collection
var CMDB_Inventory_DB *mongo.Collection
var CMDB_Pending_DB *mongo.Collection
var Core_Proposition_DB *mongo.Collection
var Core_Jobs_DB *mongo.Collection
var Core_Capability_DB *mongo.Collection
var Core_Cookbooks_DB *mongo.Collection

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
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	utils.ErrorFatal("MongoDB couldn't make initial connection to "+uri, err)

	// Ping the primary
	err = client.Ping(context.TODO(), readpref.Primary())
	utils.ErrorFatal("MongoDB couldn't ping "+uri, err)
	utils.Log("Successfully connected MongoDB to: "+uri, false)

	Results_Nmap_DB = client.Database("Results").Collection("Nmap")
	utils.Log("Successfully setup mongo nmap database collections: ", false)

	CMDB_Inventory_DB = client.Database("CMDB").Collection("Inventory")
	utils.Log("Successfully setup mongo Inventory database collections: ", false)

	CMDB_Pending_DB = client.Database("CMDB").Collection("Pending")
	utils.Log("Successfully setup mongo Pending database collections: ", false)

	Core_Proposition_DB = client.Database("Core").Collection("Proposition")
	utils.Log("Successfully setup mongo Proposition database collections: ", false)

	Core_Jobs_DB = client.Database("Core").Collection("Jobs")
	utils.Log("Successfully setup mongo Jobs database collections: ", false)

	Core_Capability_DB = client.Database("Core").Collection("Capability")
	utils.Log("Successfully setup mongo capability database collections: ", false)

	Core_Cookbooks_DB = client.Database("Core").Collection("Cookbooks")
	utils.Log("Successfully setup mongo cookbooks database collections: ", false)
}
