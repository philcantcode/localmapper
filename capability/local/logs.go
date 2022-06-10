package local

import (
	"net/http"

	"github.com/philcantcode/localmapper/system"
)

/* HTTP_JSON_GetNetworkAdapters returns all network adapters on the server */
func HTTP_JSON_GetLogs(w http.ResponseWriter, r *http.Request) {
	//TODO:: give logs

	system.Force("NOT IMPLEMENTED", true)

	//json.NewEncoder(w).Encode(logList)
}
