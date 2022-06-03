package local

import (
	"encoding/json"
	"net/http"

	"github.com/philcantcode/localmapper/utils"
)

/* HTTP_JSON_GetNetworkAdapters returns all network adapters on the server */
func HTTP_JSON_GetLogs(w http.ResponseWriter, r *http.Request) {
	logs := utils.ParseFile(utils.Configs["JSON_LOG"])
	logList := []utils.JsonLog{}

	// Skip over the last line with a newline
	for _, logEntry := range logs {

		if len(logEntry) == 0 {
			continue
		}

		jsonLog := utils.JsonLog{}
		json.Unmarshal([]byte(logEntry), &jsonLog)

		logList = append(logList, jsonLog)
	}

	json.NewEncoder(w).Encode(logList)
}
