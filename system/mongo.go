package system

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connection URI
var uri string
var client *mongo.Client
var System_Logs_DB *mongo.Collection
var Results_Nmap_DB *mongo.Collection
var CMDB_Inventory_DB *mongo.Collection
var CMDB_Pending_DB *mongo.Collection
var Core_Proposition_DB *mongo.Collection
var Core_Jobs_DB *mongo.Collection
var Core_Capability_DB *mongo.Collection
var Core_Cookbooks_DB *mongo.Collection

var MONGO_INITIALISED = false

func InitMongo() {
	var err error

	if Get("mongo-password-req") == "1" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s",
			Get("mongo-user"),
			Get("mongo-password"),
			Get("mongo-ip"),
			Get("mongo-port"))
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s",
			Get("mongo-ip"),
			Get("mongo-port"))
	}

	Log("Attempting to connect MongoDB to: "+uri, true)

	// Create a new client and connect to the server
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	Fatal("MongoDB couldn't make initial connection to "+uri, err)

	// Ping the primary
	err = client.Ping(context.TODO(), readpref.Primary())
	Fatal("Can't reach mongo at: "+uri, err)

	Log("Successfully connected MongoDB to: "+uri, false)

	Results_Nmap_DB = client.Database("Results").Collection("Nmap")
	Log("Successfully setup mongo nmap database collections: ", false)

	CMDB_Inventory_DB = client.Database("CMDB").Collection("Inventory")
	Log("Successfully setup mongo Inventory database collections: ", false)

	CMDB_Pending_DB = client.Database("CMDB").Collection("Pending")
	Log("Successfully setup mongo Pending database collections: ", false)

	Core_Proposition_DB = client.Database("Core").Collection("Proposition")
	Log("Successfully setup mongo Proposition database collections: ", false)

	Core_Jobs_DB = client.Database("Core").Collection("Jobs")
	Log("Successfully setup mongo Jobs database collections: ", false)

	Core_Capability_DB = client.Database("Core").Collection("Capability")
	Log("Successfully setup mongo capability database collections: ", false)

	Core_Cookbooks_DB = client.Database("Core").Collection("Cookbooks")
	Log("Successfully setup mongo cookbooks database collections: ", false)

	System_Logs_DB = client.Database("System").Collection("Logs")
	Log("Successfully setup mongo system logs database collections: ", true)

	MONGO_INITIALISED = true
}
