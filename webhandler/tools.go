package webhandler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/tools/nmap"
	"go.mongodb.org/mongo-driver/bson"
)

type ToolsHandler struct {
	Nmap         NmapHandler
	Searchsploit SearchsploitHandler
}

type NmapHandler struct {
}

type SearchsploitHandler struct {
}

var Tools = ToolsHandler{}
var Nmap = NmapHandler{}
var Searchsploit = SearchsploitHandler{}

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

func (sp *SearchsploitHandler) HTTP_FILE_ServeFile(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")

	installPath := system.GetConfig("external-resources-path") + "/exploitdb"

	if !strings.HasPrefix(path, installPath) {
		system.Warning("Searchsploit path error, outside of dir", true)
	}

	system.Log("Searchsploit downloading file: "+path, true)

	http.ServeFile(w, r, path)
}
