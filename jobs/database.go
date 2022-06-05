package jobs

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FILTER_JobSpec(filter bson.M, projection bson.M) []JobSpec {
	cursor, err := database.JobsDB.Find(context.Background(), filter, options.Find().SetProjection(projection))
	utils.ErrorFatal("Couldn't SELECT_JobSpec", err)
	defer cursor.Close(context.Background())

	var results []JobSpec

	for cursor.Next(context.Background()) {
		var jobSpec JobSpec

		err = cursor.Decode(&jobSpec)
		utils.ErrorFatal("Couldn't decode JobSpec", err)

		results = append(results, jobSpec)
	}

	return results
}

func INSERT_JobSpec(job JobSpec) {
	utils.Log("Attempting to INSERT_JobSpec", false)

	job.ID = primitive.NewObjectID()
	insertResult, err := database.JobsDB.InsertOne(context.Background(), job)

	utils.ErrorFatal("Couldn't INSERT_JobSpec", err)
	utils.Log(fmt.Sprintf("New Insert at: %s", insertResult), true)
}
