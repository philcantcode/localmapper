package capability

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/database"
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

func HTTP_JSON_GetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SELECT_Capability(bson.M{}, bson.M{}))
}

/* HTTP_JSON_GetCMDBCompatible returns a list of capabilities that can be
run by a particular CMDB item given it's Tags */
func HTTP_JSON_GetCMDBCompatible(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	result := []Capability{}

	entry := cmdb.SELECT_ENTRY_Inventory(bson.M{"_id": database.EncodeID(id)}, bson.M{})[0]
	caps := SELECT_Capability(bson.M{}, bson.M{})

	for _, cap := range caps {
		isCompatible, parsedCap := CMP_Entry_Capability(cap, entry)

		if isCompatible {
			result = append(result, parsedCap)
		}
	}

	json.NewEncoder(w).Encode(result)
}

func HTTP_JSON_GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	capabilities := SELECT_Capability(bson.M{"_id": database.EncodeID(id)}, bson.M{})
	json.NewEncoder(w).Encode(capabilities[0])
}

/* runCapability executes one specific capability */
func HTTP_JSON_Run(w http.ResponseWriter, r *http.Request) {
	capabilityParam := r.FormValue("capability")
	var capability Capability

	json.Unmarshal([]byte(capabilityParam), &capability)

	result := ProcessCapability(capability)

	json.NewEncoder(w).Encode(result)
}

/*
HTTP_JSON_RunCMDBCompatible takes in 2 IDs for a capability and CMDB entry
and finds any matching capabilities given the CMDB SysTags
*/
func HTTP_JSON_RunCMDBCompatible(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cmbd_id := params["cmbd_id"]
	cap_id := params["capability_id"]

	cap := SELECT_Capability(bson.M{"_id": database.EncodeID(cap_id)}, bson.M{})[0]
	entry := cmdb.SELECT_ENTRY_Inventory(bson.M{"_id": database.EncodeID(cmbd_id)}, bson.M{})[0]

	isCompatible, parsedCap := CMP_Entry_Capability(cap, entry)

	if isCompatible {
		w.Write(ProcessCapability(parsedCap))
		return
	}

	utils.ErrorContextLog("HTTP_JSON_RunCMDBCompatible was not compatible", true)
}

func CMP_Entry_Capability(capability Capability, entry cmdb.Entry) (bool, Capability) {
	var success bool

	for k, capParam := range capability.Command.Params {
		success, capability.Command.Params[k] = CMP_CapabilityParam_Entry(capParam, entry.SysTags)

		if !success {
			return false, capability
		}
	}

	return true, capability
}

/*
Determines if given a capability param {"Value": "","DataType": 1, "Default": ""}
Is there any SysTags that can fulfil the Values
*/
func CMP_CapabilityParam_Entry(capParam Param, entryTags []cmdb.EntryTag) (bool, Param) {
	// For each: {DataType.CMDB, DataType.IP}
	for _, pType := range capParam.DataType {
		// If the value is already set, move on
		if capParam.Value != "" {
			return true, capParam
		}

		// Skip empty tags that don't require input
		if pType == utils.EMPTY {
			return true, capParam
		}

		for _, eTag := range entryTags {
			// The DataType match
			if pType == eTag.DataType {
				capParam.Value = eTag.Values[len(eTag.Values)-1]
				return true, capParam
			}
		}
	}

	return false, capParam
}
