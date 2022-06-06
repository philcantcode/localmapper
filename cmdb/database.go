package cmdb

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func INSERT_CMDBItem(cmdb CMDBItem) {
	utils.Log("Attempting to INSERT_CMDB", false)

	cmdb.ID = primitive.NewObjectID()
	insertResult, err := database.CMDB_Devices_DB.InsertOne(context.Background(), cmdb)

	utils.ErrorFatal("Couldn't INSERT_CMDB", err)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)
}

func SELECT_CMDBItem(filter bson.M, projection bson.M) []CMDBItem {
	cursor, err := database.CMDB_Devices_DB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	utils.ErrorFatal("Couldn't SELECT_CMDBItem", err)
	defer cursor.Close(context.Background())

	results := []CMDBItem{}

	for cursor.Next(context.Background()) {
		var cmdb CMDBItem

		err = cursor.Decode(&cmdb)
		utils.ErrorFatal("Couldn't decode SELECT_CMDBItem", err)

		results = append(results, cmdb)
	}

	return results
}

func INSERT_VLAN(vlan Vlan) {
	utils.Log("Attempting to INSERT_VLAN", false)

	vlan.ID = primitive.NewObjectID()
	insertResult, err := database.CMDB_VLAN_DB.InsertOne(context.Background(), vlan)

	utils.ErrorFatal("Couldn't INSERT_VLAN", err)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)
}

func SELECT_Vlan(filter bson.M, projection bson.M) []Vlan {
	cursor, err := database.CMDB_VLAN_DB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	utils.ErrorFatal("Couldn't SELECT_Vlan", err)
	defer cursor.Close(context.Background())

	results := []Vlan{}

	for cursor.Next(context.Background()) {
		var vlan Vlan

		err = cursor.Decode(&vlan)
		utils.ErrorFatal("Couldn't decode SELECT_Vlan", err)

		results = append(results, vlan)
	}

	return results
}
