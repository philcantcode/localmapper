package webhandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/proposition"
	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
)

type CMDBHandler struct {
}

var CMDB = CMDBHandler{}

func (entity *CMDBHandler) HTTP_JSON_Inventory_GetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(cmdb.SELECT_ENTRY_Inventory(bson.M{}, bson.M{}))
}

func (entity *CMDBHandler) HTTP_JSON_Pending_GetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(cmdb.SELECT_ENTRY_Pending(bson.M{}, bson.M{}))
}

func (entity *CMDBHandler) HTTP_JSON_GetLocal(w http.ResponseWriter, r *http.Request) {
	local := cmdb.SELECT_ENTRY_Inventory(bson.M{"systags.label": "Identity", "systags.values": "local"}, bson.M{})

	if len(local) != 1 {
		json.NewEncoder(w).Encode(cmdb.Entity{})
		return
	}

	json.NewEncoder(w).Encode(local[0])
}

func (entity *CMDBHandler) HTTP_JSON_GetByType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	typeStr := params["type"]
	typeNum, err := strconv.Atoi(typeStr)
	system.Error("Couldn't convert CMDBType to int32: "+typeStr, err)

	json.NewEncoder(w).Encode(cmdb.SELECT_ENTRY_Inventory(bson.M{"cmdbtype": typeNum}, bson.M{}))
}

/* HTTP_JSON_GetByID returns a CMDB element based on the {id} */
func (entity *CMDBHandler) HTTP_JSON_GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	entities := cmdb.SELECT_Entities_Joined(bson.M{"_id": system.EncodeID(id)}, bson.M{})

	if len(entities) != 1 {
		system.Warning("Incorrect number of Entities returned by ID in Webhandler", true)
	}

	json.NewEncoder(w).Encode(entities)
}

func (entity *CMDBHandler) HTTP_NONE_ApprovePending(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("ID")

	pending := cmdb.SELECT_ENTRY_Pending(bson.M{"_id": system.EncodeID(id)}, bson.M{})

	if len(pending) != 1 {
		system.Warning("Incorrect number of Entities returned by ID in ApprovePending", true)
	}

	pending[0].ApprovePending()
}

func (entity *CMDBHandler) HTTP_NONE_DenyPending(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("ID")

	pending := cmdb.SELECT_ENTRY_Pending(bson.M{"_id": system.EncodeID(id)}, bson.M{})

	if len(pending) != 1 {
		system.Warning("Incorrect number of Entities returned by ID in DenyPending", true)
	}

	pending[0].DELETE_ENTRY_Pending()
}

func (entity *CMDBHandler) HTTP_NONE_DenyAll(w http.ResponseWriter, r *http.Request) {
	pending := cmdb.SELECT_ENTRY_Pending(bson.M{}, bson.M{})

	for _, entity := range pending {
		entity.DELETE_ENTRY_Pending()
	}
}

func (entity *CMDBHandler) HTTP_JSON_GetStats(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	results := cmdb.SELECT_Entities_Joined(bson.M{"_id": system.EncodeID(id)}, bson.M{})

	if len(results) != 1 {
		system.Warning("Incorrect number results returned getting HTTP_JSON_GetStats", true)
	}

	json.NewEncoder(w).Encode(results[0].GetStats())
}

/*
	HTTP_JSON_Restore restores the system settings and databases
	to factory defaults.
*/
func (entity *CMDBHandler) HTTP_JSON_Restore(w http.ResponseWriter, r *http.Request) {
	cmdb.Restore()
	w.Write([]byte("200/Done"))
}

//TODO: Add validation
func (entity *CMDBHandler) HTTP_NONE_NewVLAN(w http.ResponseWriter, r *http.Request) {
	label := r.FormValue("Label")
	desc := r.FormValue("Desc")
	highIP := r.FormValue("HighIP")
	lowIP := r.FormValue("LowIP")
	cmdbType := r.FormValue("CMDBType")

	cmdbTypeInt, err := strconv.Atoi(cmdbType)
	system.Error("Couldn't convert CMDBType to int", err)

	cmdb.CreateNewVLAN(label, desc, lowIP, highIP, cmdbTypeInt)

	w.Write([]byte("200/Success"))
}

func (entity *CMDBHandler) HTTP_None_ResolveLocalIdentityProposition(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("ID")
	value := r.PostFormValue("Value")

	valueIdx, err := strconv.Atoi(value)
	system.Error("Couldn't convert proposition value ID to int", err)

	for i, prop := range proposition.Propositions {
		if prop.ID == id {
			cmdb.SetLocalIdentityEntry(proposition.Propositions[i].Predicates[valueIdx].Value)
			proposition.Pop(valueIdx)
			break
		}
	}
}

func (entity *CMDBHandler) HTTP_None_ResolveIPConflict(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("ID")
	value := r.PostFormValue("Value")

	valueIdx, err := strconv.Atoi(value)
	system.Error("Couldn't convert proposition value ID to int", err)

	for _, prop := range proposition.Propositions {
		if prop.ID == id {
			cmdb.ResolveIPConflict(cmdb.ConflictActions(prop.Predicates[valueIdx].Value), prop.GetEvidenceValue("Conflict IP"))
			proposition.Pop(valueIdx)
			break
		}
	}
}
