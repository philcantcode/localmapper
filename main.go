package main

import (
	"github.com/philcantcode/localmapper/adapters/cli"
	"github.com/philcantcode/localmapper/adapters/web"
	"github.com/philcantcode/localmapper/application/database"
	"github.com/philcantcode/localmapper/utils"
)

func main() {
	utils.LoadGlobalConfigs()

	database.InitSqlite()
	database.InitMongo()

	//go network.PingSweepVlans()

	cli.InitCLI()
	web.InitServer()
}
