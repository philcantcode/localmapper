package system

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// func ExecuteCommand(id int) {
// 	Log(fmt.Sprintf("Executing System Command: %d", id), true)

// 	switch id {
// 	case 1:
// 		Core_Capability_DB.Drop(context.Background())
// 	case 2:
// 		CMDB_Pending_DB.Drop(context.Background())
// 	case 3:
// 		CMDB_Inventory_DB.Drop(context.Background())
// 	case 4:
// 		Core_Proposition_DB.Drop(context.Background())
// 	case 5:
// 		Results_Nmap_DB.Drop(context.Background())
// 	case 6:
// 		proposition.FirstTimeSetup()
// 	case 7:
// 		capability.InsertDefaultCapabilities()
// 	case 8:
// 		Core_Capability_DB.Drop(context.Background())
// 		CMDB_Pending_DB.Drop(context.Background())
// 		CMDB_Inventory_DB.Drop(context.Background())
// 		Core_Proposition_DB.Drop(context.Background())
// 		Core_Cookbooks_DB.Drop(context.Background())
// 		Results_Nmap_DB.Drop(context.Background())
// 		proposition.FirstTimeSetup()
// 		capability.InsertDefaultCapabilities()
// 		cookbook.InsertDefaultCookbooks()
// 	case 9:
// 		Core_Cookbooks_DB.Drop(context.Background())
// 	}
// }

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
