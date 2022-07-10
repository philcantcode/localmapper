package system

import (
	"context"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var SystemJobs = []ScheduledJob{}

func GetConfig(key string) string {
	for _, k := range SELECT_Settings_All() {
		if k.Key == key {
			return k.Value
		}
	}

	Warning("Couldn't retrieve key: "+key, true)
	return ""
}

func GetInt(key string) int {
	for _, k := range SELECT_Settings_All() {
		if k.Key == key {
			intKey, err := strconv.Atoi(k.Value)
			Error("Couldn't convert key to int", err)

			return intKey
		}
	}

	Warning("Couldn't retrieve key: "+key, true)
	return -1
}

func EncodeID(id string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(id)
	Error("Couldn't convert string ID to mongo ID Object", err)

	return objID
}

func Restore() {
	Log("Restoring system settings to factory defaults", true)

	System_Logs_DB.Drop(context.Background()) // Drop the logs table
	DELETE_Settings_All()                     // Delete all settings

	Init() // Perform first time setup
}
