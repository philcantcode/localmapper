package nmap

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type KeyValue struct {
	Key   string
	Value string
}

func HTTP_SELECT_Logs(w http.ResponseWriter, r *http.Request) {
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

	results := SELECT(filterBson, projBson)

	fmt.Printf("RETUNRING LEN: %d\n", len(results))

	json.NewEncoder(w).Encode(results)
}
