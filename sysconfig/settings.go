package sysconfig

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/capability"
	cookbook "github.com/philcantcode/localmapper/capability/cookbooks"
	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/proposition"
	"github.com/philcantcode/localmapper/utils"
)

func ExecuteCommand(id int) {
	utils.Log(fmt.Sprintf("Executing SysConfig Command: %d", id), true)

	switch id {
	case 1:
		database.Core_Capability_DB.Drop(context.Background())
	case 2:
		database.CMDB_Pending_DB.Drop(context.Background())
	case 3:
		database.CMDB_Inventory_DB.Drop(context.Background())
	case 4:
		database.Core_Proposition_DB.Drop(context.Background())
	case 5:
		database.Results_Nmap_DB.Drop(context.Background())
	case 6:
		proposition.SetupJobs()
	case 7:
		capability.InsertDefaultCapabilities()
	case 8:
		database.Core_Capability_DB.Drop(context.Background())
		database.CMDB_Pending_DB.Drop(context.Background())
		database.CMDB_Inventory_DB.Drop(context.Background())
		database.Core_Proposition_DB.Drop(context.Background())
		database.Core_Cookbooks_DB.Drop(context.Background())
		database.Results_Nmap_DB.Drop(context.Background())
		proposition.SetupJobs()
		capability.InsertDefaultCapabilities()
		cookbook.InsertDefaultCookbooks()
	case 9:
		database.Core_Cookbooks_DB.Drop(context.Background())
	}
}
