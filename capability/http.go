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

	entry := cmdb.SELECT_ENTRY(bson.M{"_id": database.EncodeID(id)}, bson.M{})[0]
	caps := SELECT_Capability(bson.M{}, bson.M{})

	entryDataTypes := []utils.DataType{}

	for _, v := range entry.SysTags {
		entryDataTypes = append(entryDataTypes, v.DataType)
	}

outer:
	for _, cap := range caps {
		for _, param := range cap.Command.Params {
			switch param.DataType {
			case utils.EMPTY:
				continue
			default:
				if !utils.DTArrayContains(param.DataType, entryDataTypes) && param.Default == "" {
					continue outer
				}
			}
		}

		result = append(result, cap)
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

func HTTP_JSON_RunCMDBCompatible(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cmbd_id := params["cmbd_id"]
	cap_id := params["capability_id"]

	cap := SELECT_Capability(bson.M{"_id": database.EncodeID(cap_id)}, bson.M{})[0]
	entry := cmdb.SELECT_ENTRY(bson.M{"_id": database.EncodeID(cmbd_id)}, bson.M{})[0]

	for k, param := range cap.Command.Params {
		for _, tag := range entry.SysTags {
			if tag.DataType == param.DataType && param.Default == "" {
				cap.Command.Params[k].Value = tag.Values[0]
			}
		}
	}

	result := ProcessCapability(cap)
	w.Write(result)
}
