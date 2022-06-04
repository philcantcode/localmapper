package main

import (
	"github.com/philcantcode/localmapper/api"
	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/propositions"
	"github.com/philcantcode/localmapper/utils"
)

func main() {
	utils.LoadGlobalConfigs()

	database.InitSqlite()
	database.InitMongo()

	propositions.SetupJobs()

	api.InitServer()
}
