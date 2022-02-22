package main

import (
	"fmt"

	"github.com/philcantcode/localmapper/adapters/cli"
	"github.com/philcantcode/localmapper/adapters/network"
	"github.com/philcantcode/localmapper/adapters/web"
	"github.com/philcantcode/localmapper/application/database"
	"github.com/philcantcode/localmapper/utils"
)

func main() {
	utils.LoadGlobalConfigs()

	database.InitSqlite()
	database.InitMongo()

	fmt.Println(utils.PrettyPrint(network.ListAllAddresses()))

	cli.InitCLI()
	web.InitServer()
}
