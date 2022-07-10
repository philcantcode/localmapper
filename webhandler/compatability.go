package webhandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/compatability"
	"github.com/philcantcode/localmapper/system"
)

type CompatabilityHandler struct {
}

var Compatability = CompatabilityHandler{}

/*
	HTTP_JSON_GetCMDBCompatible returns a list of capabilities that can be
	run by a particular CMDB item given it's Tags
*/
func (compat *CompatabilityHandler) HTTP_JSON_GET_Capability_ByEntityID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["entityID"]

	json.NewEncoder(w).Encode(compatability.FindCompatibleCapabilitiesForEntity(id))
}

/*
	HTTP_NONE_RunCapability takes in 2 IDs for a capability and CMDB entry
	and finds any matching capabilities given the CMDB SysTags
*/
func (compat *CompatabilityHandler) HTTP_NONE_RunCapability(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	r.ParseForm()
	cmdb_id := params["entityID"]
	cap_id := params["capabilityID"]
	options := make(map[string]int)

	/*
		Parse the POST params which are just options KEY VALUE pairs
		The key is the flag, the value[0] is the index of the desired option.
	*/
	for key, val := range r.PostForm {
		optID, err := strconv.Atoi(val[0])
		system.Error("Could not convert from optionID to int", err)

		options[key] = optID
	}

	compatability.LaunchCapabilityForEntity(cap_id, cmdb_id, options)
}
