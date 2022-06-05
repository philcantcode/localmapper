package jobs

import (
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobSpec struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string
	Desc         string
	DataType     utils.DataType // Value type
	Targets      []string       // Targets (e.g., IPs)
	Capabilities []int          // Capabilities to be performed on the data
}

func TEST_CreateJobSpec1() {
	jobSpec := JobSpec{Name: "Ping Sweep", Desc: "Perform a Ping Sweep of an IP range.", DataType: utils.IPRange, Capabilities: []int{2}, Targets: []string{}}
	INSERT_JobSpec(jobSpec)
}
