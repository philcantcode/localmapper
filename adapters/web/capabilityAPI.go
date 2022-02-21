package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/philcantcode/localmapper/adapters/blueprint"
	"github.com/philcantcode/localmapper/application/database"
	"github.com/philcantcode/localmapper/application/nmap"
	"github.com/philcantcode/localmapper/utils"
)

/* updateCapability takes in a single capability (JSON object)
   and updates it via the ID */
func updateCapability(w http.ResponseWriter, r *http.Request) {
	capabilityParam := r.FormValue("capability")
	var capability blueprint.Capability

	err := json.Unmarshal([]byte(capabilityParam), &capability)
	utils.ErrorFatal("Error converting capability (json string) > capability (struct)", err)

	database.UpdateCapability(capability)
	w.WriteHeader(200)
}

/* getCapabilities returns all capabilities as JSON,
   if an ID is specified, it returns that capability,
   otherwise all are returned */
func getCapabilities(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	capabilities := database.SelectAllCapabilities()

	if id == "" {
		json.NewEncoder(w).Encode(capabilities)
		return
	}

	capabilityID, err := strconv.Atoi(id)
	utils.ErrorFatal("Couldn't convert ID in GetCapabilities", err)

	for _, capability := range capabilities {
		if capability.ID == capabilityID {
			json.NewEncoder(w).Encode(capability)
			return
		}
	}
}

/* runCapability executes one specific capability */
func runCapability(w http.ResponseWriter, r *http.Request) {
	capabilityParam := r.FormValue("capability")
	var capability blueprint.Capability

	json.Unmarshal([]byte(capabilityParam), &capability)

	switch capability.Type {
	case "nmap":
		nmapRun := nmap.RunNmapCommand(capability)
		database.InsertNetworkNmap(nmapRun)
		utils.PrintLog(utils.PrettyPrintToStr(nmapRun))
		json.NewEncoder(w).Encode(nmapRun)
		return
	default:
		utils.ErrorForceFatal(fmt.Sprintf("No capability type to run in adapters.api.RunCapability: %v", capability))
	}
}
