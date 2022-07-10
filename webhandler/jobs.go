package webhandler

import (
	"encoding/json"
	"net/http"

	"github.com/philcantcode/localmapper/capability"
)

type JobsHandler struct {
}

var Jobs = JobsHandler{}

func (job *JobsHandler) HTTP_JSON_GetStats(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(capability.GetJobStats())
}
