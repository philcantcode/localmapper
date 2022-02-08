package api

import (
	"encoding/json"
	"net/http"

	"github.com/philcantcode/localmapper/discovery"
	"github.com/philcantcode/localmapper/utils"
)

func init() {

}

func NmapPingScan(w http.ResponseWriter, r *http.Request) {

	target := r.FormValue("target")

	if target == "" {
		utils.ErrorHandleLog("No target specified for NmapPingScan", true)
		return
	}

	xml := discovery.PingScan(target)

	json.NewEncoder(w).Encode(xml.Hosts)
}

func NmapOSDetectionScan(w http.ResponseWriter, r *http.Request) {

	target := r.FormValue("target")

	if target == "" {
		utils.ErrorHandleLog("No target specified for NmapPingScan", true)
		return
	}

	xml := discovery.OSDetectionScan(target)

	json.NewEncoder(w).Encode(xml)
}
