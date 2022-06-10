package main

import (
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/philcantcode/localmapper/capability"
	"github.com/philcantcode/localmapper/capability/cookbook"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/proposition"
	"github.com/philcantcode/localmapper/system"
)

func main() {
	// Load settings & config database
	system.InitSqlite()
	system.FirstTimeSetup()

	// Load application database
	system.InitMongo()

	// Load all initial setup jobs here
	proposition.FirstTimeSetup()
	capability.FirstTimeSetup()
	cookbook.FirstTimeSetup()
	cmdb.FirstTimeSetup()

	// Initialise the web API
	initServer()
}

func omiTest() {
	// init COM, oh yeah
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unknown, _ := oleutil.CreateObject("WbemScripting.SWbemLocator")
	defer unknown.Release()

	wmi, _ := unknown.QueryInterface(ole.IID_IDispatch)
	defer wmi.Release()

	// service is a SWbemServices
	serviceRaw, _ := oleutil.CallMethod(wmi, "ConnectServer")
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// result is a SWBemObjectSet
	resultRaw, _ := oleutil.CallMethod(service, "ExecQuery", "SELECT * FROM Win32_Process")
	result := resultRaw.ToIDispatch()
	defer result.Release()

	countVar, _ := oleutil.GetProperty(result, "Count")
	count := int(countVar.Val)

	for i := 0; i < count; i++ {
		// item is a SWbemObject, but really a Win32_Process
		itemRaw, _ := oleutil.CallMethod(result, "ItemIndex", i)
		item := itemRaw.ToIDispatch()
		defer item.Release()

		asString, _ := oleutil.GetProperty(item, "Name")

		println(asString.ToString())
	}
}

func ExecuteCommand(id int) {
	system.Log(fmt.Sprintf("Executing System Command: %d", id), true)

	// switch id {
	// case 1:
	// 	system.Core_Capability_DB.Drop(context.Background())
	// case 2:
	// 	system.CMDB_Pending_DB.Drop(context.Background())
	// case 3:
	// 	system.CMDB_Inventory_DB.Drop(context.Background())
	// case 4:
	// 	system.Core_Proposition_DB.Drop(context.Background())
	// case 5:
	// 	system.Results_Nmap_DB.Drop(context.Background())
	// case 6:
	// 	proposition.FirstTimeSetup()
	// case 7:
	// 	capability.InsertDefaultCapabilities()
	// case 8:
	// 	system.Core_Capability_DB.Drop(context.Background())
	// 	system.CMDB_Pending_DB.Drop(context.Background())
	// 	system.CMDB_Inventory_DB.Drop(context.Background())
	// 	system.Core_Proposition_DB.Drop(context.Background())
	// 	system.Core_Cookbooks_DB.Drop(context.Background())
	// 	system.Results_Nmap_DB.Drop(context.Background())
	// 	proposition.FirstTimeSetup()
	// 	capability.InsertDefaultCapabilities()
	// 	cookbook.InsertDefaultCookbooks()
	// case 9:
	// 	system.Core_Cookbooks_DB.Drop(context.Background())
	// }
}
