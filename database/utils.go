package database

import (
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func EncodeID(id string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(id)
	utils.ErrorLog("Couldn't convert string ID to mongo ID Object", err, false)

	return objID
}
