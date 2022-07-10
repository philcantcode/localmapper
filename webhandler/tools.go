package webhandler

import (
	"encoding/json"
	"net/http"

	"github.com/philcantcode/localmapper/tools/nmap"
	"go.mongodb.org/mongo-driver/bson"
)

type ToolsHandler struct {
}

var Tools = ToolsHandler{}

func (tool *ToolsHandler) HTTP_JSON_SELECT_Logs(w http.ResponseWriter, r *http.Request) {
	filter := r.FormValue("filter")         // JSON {K:V}
	projection := r.FormValue("projection") // JSON {K:V}

	var filters = []map[string]string{}
	json.Unmarshal([]byte(filter), &filters)

	var projections = []map[string]string{}
	json.Unmarshal([]byte(projection), &projections)

	filterBson := bson.M{}
	projBson := bson.M{}

	for _, filter := range filters {
		for key, val := range filter {
			filterBson[key] = val
		}
	}

	for _, projection := range projections {
		for key, val := range projection {
			projBson[key] = val
		}
	}

	results := nmap.SELECT(filterBson, projBson)
	json.NewEncoder(w).Encode(results)
}
