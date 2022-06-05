package capability

import (
	"encoding/json"
	"net/http"

	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

/* updateCapability takes in a single capability (JSON object)
   and updates it via the ID */
func HTTP_JSON_Update(w http.ResponseWriter, r *http.Request) {
	capabilityParam := r.FormValue("capability")
	var capability Capability

	err := json.Unmarshal([]byte(capabilityParam), &capability)
	utils.ErrorFatal("Error converting capability (json string) > capability (struct)", err)

	UPDATE_Capability(capability)
	w.WriteHeader(200)
}

/* getCapabilities returns all capabilities as JSON,
   if an ID is specified, it returns that capability,
   otherwise all are returned */
func HTTP_JSON_Get(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	capabilities := SELECT_Capability(bson.M{}, bson.M{})

	if id == "" {
		json.NewEncoder(w).Encode(capabilities)
		return
	}

	for _, capability := range capabilities {
		if capability.ID.String() == id {
			json.NewEncoder(w).Encode(capability)
			return
		}
	}
}

/* runCapability executes one specific capability */
func HTTP_JSON_Run(w http.ResponseWriter, r *http.Request) {
	capabilityParam := r.FormValue("capability")
	var capability Capability

	json.Unmarshal([]byte(capabilityParam), &capability)

	result := ProcessCapability(capability)

	json.NewEncoder(w).Encode(result)
}
