package webhandler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/capability"
	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
)

type CapabilityHandler struct {
}

var Capability = CapabilityHandler{}

/*
	INPUT: None

	OUTPUT: Array(Capability)
*/
func (cap *CapabilityHandler) HTTP_JSON_GetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(capability.SELECT_Capability(bson.M{}, bson.M{}))
}

/*
	INPUT: String(Capability)

	PROCESS: String(Capability) > JSON(Unmarshal) > Capability > Update

	OUTPUT: None
*/
func (cap *CapabilityHandler) HTTP_JSON_Update(w http.ResponseWriter, r *http.Request) {
	capabilityParam := r.FormValue("capability")
	var capability capability.Capability

	err := json.Unmarshal([]byte(capabilityParam), &capability)
	system.Fatal("Error converting capability (json string) > capability (struct)", err)

	capability.UPDATE_Capability()
	w.WriteHeader(200)
}

/*
	INPUT: Int(ID)

	OUTPUT: Single Capability
*/
func (cap *CapabilityHandler) HTTP_JSON_GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	capabilities := capability.SELECT_Capability(bson.M{"_id": system.EncodeID(id)}, bson.M{})

	if len(capabilities) > 1 {
		system.Warning("Too many capabilities returned in HTTP_JSON_GetByID", true)
	}

	json.NewEncoder(w).Encode(capabilities[0])
}

/*
	INPUT: String(Capability)

	OUTPUT: None
*/
func (cap *CapabilityHandler) HTTP_NONE_Run(w http.ResponseWriter, r *http.Request) {
	capabilityParam := r.FormValue("capability")
	var capability capability.Capability

	json.Unmarshal([]byte(capabilityParam), &capability)
	capability.LaunchLifecycleManager()

	json.NewEncoder(w).Encode("200/Done")
}

/*

 */
func (cap *CapabilityHandler) HTTP_JSON_BLANK_Capability(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(capability.Capability{})
}

func (cap *CapabilityHandler) HTTP_JSON_BLANK_Command(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(capability.Command{})
}

func (cap *CapabilityHandler) HTTP_JSON_BLANK_Param(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(capability.Param{})
}

/*
	HTTP_JSON_Restore restores the system settings and databases
	to factory defaults.
*/
func (cap *CapabilityHandler) HTTP_NONE_Restore(w http.ResponseWriter, r *http.Request) {

	capability.Restore()
	w.Write([]byte("200/Done"))
}
