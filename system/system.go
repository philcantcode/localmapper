package system

import (
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

	Force("Couldn't retrieve key: "+key, true)
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

	Force("Couldn't retrieve key: "+key, true)
	return -1
}

func EncodeID(id string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(id)
	Error("Couldn't convert string ID to mongo ID Object", err)

	return objID
}
