package database

import (
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ID_TO_Obj(id string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(id)
	utils.ErrorFatal("Couldn't convert string ID to mongo ID Object", err)

	return objID
}
