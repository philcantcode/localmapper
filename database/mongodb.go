package database

import (
	"context"
	"fmt"
	"log"

	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection URI
var uri string
var client *mongo.Client
var collection *mongo.Collection

func MongoConnect() {
	var err error
	uri = fmt.Sprintf("mongodb://%s:%s@%s:%s",
		utils.Configs["MONGO_USER"],
		utils.Configs["MONGO_PASSWORD"],
		utils.Configs["MONGO_IP"],
		utils.Configs["MONGO_PORT"])

	// Create a new client and connect to the server
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// // Ping the primary
	// if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
	// 	panic(err)
	// }

	collection = client.Database("devices").Collection("nmap")

	fmt.Println("Successfully connected and pinged.")
}

func InsertMong(v interface{}) {
	insertResult, err := collection.InsertOne(context.TODO(), v)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}
