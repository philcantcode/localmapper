package system

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Get(key string) string {
	for _, k := range SELECT_Settings_All() {
		if k.Key == key {
			return k.Value
		}
	}

	Force("Couldn't retrieve key: "+key, true)
	return ""
}

func EncodeID(id string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(id)
	Error("Couldn't convert string ID to mongo ID Object", err)

	return objID
}
