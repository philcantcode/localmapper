package main

import (
	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/proposition"
	"github.com/philcantcode/localmapper/utils"
)

func main() {
	utils.LoadGlobalConfigs()

	database.InitSqlite()
	database.InitMongo()

	proposition.SetupJobs()

	// usrTags := []cmdb.EntryTag{}
	// usrTags = append(usrTags, cmdb.EntryTag{Label: "LowIP", Values: []string{"192.168.0.0"}})
	// usrTags = append(usrTags, cmdb.EntryTag{Label: "HighIP", Values: []string{"192.168.255.255"}})
	// sysTags := []cmdb.EntryTag{}

	// vlan := cmdb.Entry{Label: "Private Range 3", Desc: "Default VLAN", OSILayer: 2, CMDBType: cmdb.VLAN, UsrTags: usrTags, SysTags: sysTags}
	// cmdb.INSERT_ENTRY(vlan)

	initServer()
}
