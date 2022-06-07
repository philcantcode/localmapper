package main

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/proposition"
	"github.com/philcantcode/localmapper/utils"
)

func main() {
	utils.LoadGlobalConfigs()

	database.InitSqlite()
	database.InitMongo()

	proposition.SetupJobs()
	//capability.TEST_GENERATE_CAPABILITIES()
	//omiTest()

	// entry := cmdb.SELECT_ENTRY_Inventory(bson.M{"_id": database.EncodeID("629f6d27f33e5892933865a4")}, bson.M{})
	// cap := capability.SELECT_Capability(bson.M{"_id": database.EncodeID("629f808b3fb3d4ad934eef87")}, bson.M{})

	// success, res := capability.CMP_Entry_Capability(cap[0], entry[0])

	// fmt.Println(success)

	// fmt.Printf("%+v\n", res)

	// os.Exit(0)

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
