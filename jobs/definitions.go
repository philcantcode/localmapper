package jobs

import (
	"encoding/json"
	"fmt"

	"github.com/philcantcode/localmapper/utils"
)

type JobSpec struct {
	JobID        int
	Name         string
	Description  string
	DataType     utils.DataType // Value type
	Targets      []string       // Targets (e.g., IPs)
	Capabilities []int          // Capabilities to be performed on the data
}

func TEST_JobSpecCreator() {
	utils.Log("TEST_JobSpecCreator Ran", true)
	jobSpec := JobSpec{JobID: 1, Name: "Ping Sweep", Description: "Perform a Ping Sweep of an IP range.", DataType: utils.IPRange, Capabilities: []int{2}, Targets: []string{}}

	jobString, _ := json.Marshal(jobSpec)

	fmt.Printf("%+v\n\n%s\n\n", jobSpec, jobString)
	PreProcess(jobSpec)
}
