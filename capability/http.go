package capability

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
)

/* updateCapability takes in a single capability (JSON object)
   and updates it via the ID */
func HTTP_JSON_Update(w http.ResponseWriter, r *http.Request) {
	capabilityParam := r.FormValue("capability")
	var capability Capability

	err := json.Unmarshal([]byte(capabilityParam), &capability)
	system.Fatal("Error converting capability (json string) > capability (struct)", err)

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
	entries := []cmdb.Entity{}

	entries = append(entries, cmdb.SELECT_ENTRY_Inventory(bson.M{"_id": system.EncodeID(id)}, bson.M{})...)
	entries = append(entries, cmdb.SELECT_ENTRY_Pending(bson.M{"_id": system.EncodeID(id)}, bson.M{})...)

	if len(entries) != 1 {
		system.Warning("Too many results returned in HTTP_JSON_GetCMDBCompatible", true)
		return
	}

	caps := SELECT_Capability(bson.M{}, bson.M{})

	for _, cap := range caps {
		isCompatible, parsedCap := cap.CheckCompatability(entries[0])

		if isCompatible {
			result = append(result, parsedCap)
		}
	}

	json.NewEncoder(w).Encode(result)
}

func HTTP_JSON_GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	capabilities := SELECT_Capability(bson.M{"_id": system.EncodeID(id)}, bson.M{})
	json.NewEncoder(w).Encode(capabilities[0])
}

/* runCapability executes one specific capability */
func HTTP_JSON_Run(w http.ResponseWriter, r *http.Request) {
	capabilityParam := r.FormValue("capability")
	var capability Capability

	json.Unmarshal([]byte(capabilityParam), &capability)

	capability.QueueCapability()

	json.NewEncoder(w).Encode("200/Done")
}

/*
HTTP_JSON_RunCMDBCompatible takes in 2 IDs for a capability and CMDB entry
and finds any matching capabilities given the CMDB SysTags
*/
func HTTP_JSON_RunCMDBCompatible(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	r.ParseForm()
	cmdb_id := params["cmbd_id"]
	cap_id := params["capability_id"]

	cap := SELECT_Capability(bson.M{"_id": system.EncodeID(cap_id)}, bson.M{})[0]
	entries := []cmdb.Entity{}

	entries = append(entries, cmdb.SELECT_ENTRY_Joined(bson.M{"_id": system.EncodeID(cmdb_id)}, bson.M{})...)

	if len(entries) != 1 {
		system.Warning("Too many results returned in HTTP_JSON_GetCMDBCompatible", true)
		return
	}

	/*
		Parse the POST params which are just options KEY VALUE pairs
		The key is the flag, the value[0] is the index of the desired option.
	*/
	for key, val := range r.PostForm {
		optID, err := strconv.Atoi(val[0])
		system.Error("Could not convert Option ID to int", err)

		fmt.Println(val)

		for idx, param := range cap.Command.Params {
			if param.Flag == key {
				cap.Command.Params[idx].Value = param.Options[optID].Value
			}
		}
	}

	isCompatible, parsedCap := cap.CheckCompatability(entries[0])

	if isCompatible {
		parsedCap.QueueCapability()
		w.Write([]byte("200/Done"))
		return
	}

	system.Warning("HTTP_JSON_RunCMDBCompatible was not compatible", true)
}

/*
	HTTP_JSON_Restore restores the system settings and databases
	to factory defaults.
*/
func HTTP_JSON_Restore(w http.ResponseWriter, r *http.Request) {
	system.Log("Restoring capabilities to factory defaults", true)

	DELETE_Capability(bson.M{})
	system.Core_Capability_DB.Drop(context.Background()) // Drop capabilities

	FirstTimeSetup() // Restore capabilities

	w.Write([]byte("200/Done"))
}

func HTTP_JSON_GetNew(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Capability{})
}

func HTTP_JSON_GetNew_Command(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Command{})
}

func HTTP_JSON_GetNew_Param(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Param{})
}
