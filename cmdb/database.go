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
	insertResult, err := database.CmdbDB.InsertOne(context.Background(), cmdb)

	utils.ErrorFatal("Couldn't INSERT_CMDB", err)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), false)
}

func SELECT_CMDBItem(filter bson.M, projection bson.M) []CMDBItem {
	cursor, err := database.CmdbDB.Find(context.Background(), filter, options.Find().SetProjection(projection))
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

// Descending Order
func SelectAllVlans() []Vlan {
	utils.Log("SelectAllVlans from Vlans Db (sqlite)", false)
	stmt, err := database.Con.Prepare("SELECT `id`, `name`, `description`, `highIP`, `lowIP`, `tags` FROM `Vlans` ORDER BY `id` DESC")
	utils.ErrorLog("Couldn't select all from Vlans", err, true)

	rows, err := stmt.Query()
	utils.ErrorLog("Couldn't recieve rows from SelectAllVlans", err, true)
	defer rows.Close()

	vlans := []Vlan{}

	for rows.Next() {
		vlan := Vlan{}

		rows.Scan(&vlan.ID, &vlan.Name, &vlan.Description, &vlan.HighIP, &vlan.LowIP, &vlan.Tags)

		vlans = append(vlans, vlan)
	}

	return vlans
}
