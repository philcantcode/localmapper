package main

import (
	"fmt"

	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/proposition"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	utils.LoadGlobalConfigs()

	database.InitSqlite()
	database.InitMongo()

	proposition.SetupJobs()

	for _, p := range proposition.SELECT_Propositions(bson.M{}, bson.M{}) {
		fmt.Printf("%+v\n\n", p)
	}

	initServer()
}
