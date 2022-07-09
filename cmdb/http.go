package cmdb

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/proposition"
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func HTTP_JSON_GetLocal(w http.ResponseWriter, r *http.Request) {

	entries := SELECT_ENTRY_Inventory(bson.M{}, bson.M{})

	for _, entry := range entries {
		sysTag, exists, _ := entry.FindSysTag("Identity")
		if exists && utils.ArrayContains("local", sysTag.Values) {
			json.NewEncoder(w).Encode(entry)
			return
		}
	}

	json.NewEncoder(w).Encode(Entity{})
}

func HTTP_JSON_Inventory_GetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SELECT_ENTRY_Inventory(bson.M{}, bson.M{}))
}

func HTTP_JSON_Pending_GetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SELECT_ENTRY_Pending(bson.M{}, bson.M{}))
}

func HTTP_JSON_GetByType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	typeStr := params["type"]
	typeNum, err := strconv.Atoi(typeStr)
	system.Error("Couldn't convert CMDBType to int32: "+typeStr, err)

	json.NewEncoder(w).Encode(SELECT_ENTRY_Inventory(bson.M{"cmdbtype": typeNum}, bson.M{}))
}

/* HTTP_JSON_GetByID returns a CMDB element based on the {id} */
func HTTP_JSON_GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	device := []Entity{}

	system.Log("HTTP request made for: "+id, false)

	if id == "" {
		json.NewEncoder(w).Encode(SELECT_ENTRY_Inventory(bson.M{}, bson.M{}))
		return
	}

	// The inventory ID was a device
	device = append(device, SELECT_ENTRY_Inventory(bson.M{"_id": system.EncodeID(id)}, bson.M{})...)
	device = append(device, SELECT_ENTRY_Pending(bson.M{"_id": system.EncodeID(id)}, bson.M{})...)

	if len(device) > 0 {
		json.NewEncoder(w).Encode(device)
		return
	}

	// Not found return array of empty values
	json.NewEncoder(w).Encode([]Entity{})
}

//TODO: Add validation
func HTTP_INSERT_Pending_Vlan(w http.ResponseWriter, r *http.Request) {
	label := r.FormValue("Label")
	desc := r.FormValue("Desc")
	highIP := r.FormValue("HighIP")
	lowIP := r.FormValue("LowIP")
	cmdbType := r.FormValue("CMDBType")

	cmdbTypeInt, err := strconv.Atoi(cmdbType)
	system.Error("Couldn't convert CMDBType to int", err)

	if !utils.ValidateIP(lowIP) {
		system.Warning("LowIP not valid in cmdb.HTTP_INSERT_Pending", true)
		return
	}

	if !utils.ValidateIP(highIP) {
		system.Warning("HighIP not valid in cmdb.HTTP_INSERT_Pending", true)
		return
	}

	if !utils.ValidateString(label) {
		system.Warning("Label not valid in cmdb.HTTP_INSERT_Pending", true)
		return
	}

	if !utils.ValidateString(desc) {
		system.Warning("Desc not valid in cmdb.HTTP_INSERT_Pending", true)
		return
	}

	cidrArr, err := utils.IPv4RangeToCIDRRange(lowIP, highIP)
	system.Error("Couldn't create CIDR", err)

	highIpTag := EntityTag{Label: "LowIP", DataType: system.DataType_IP, Values: []string{lowIP}}
	lowIpTag := EntityTag{Label: "HighIP", DataType: system.DataType_IP, Values: []string{highIP}}
	cidr := EntityTag{Label: "CIDR", DataType: system.DataType_CIDR, Values: cidrArr}
	entry := Entity{Label: label, Description: desc, OSILayer: 2, CMDBType: CMDBType(cmdbTypeInt), DateSeen: []string{utils.GetDateTime().DateTime}, SysTags: []EntityTag{lowIpTag, highIpTag, cidr}, UsrTags: []EntityTag{}}

	entry.InsertPending()

	w.Write([]byte("200/Success"))
}

func HTTP_Pending_Approve(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("ID")

	pending := SELECT_ENTRY_Pending(bson.M{"_id": system.EncodeID(id)}, bson.M{})[0]
	pending.SysTags = append(pending.SysTags, EntityTag{Label: "Verified", DataType: system.DataType_BOOL, Values: []string{"1"}})

	pending.InsertInventory()
	DELETE_ENTRY_Pending(pending)
}

func HTTP_Pending_Deny(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("ID")

	pending := SELECT_ENTRY_Pending(bson.M{"_id": system.EncodeID(id)}, bson.M{})[0]

	DELETE_ENTRY_Pending(pending)
}

func HTTP_Pending_DenyAll(w http.ResponseWriter, r *http.Request) {
	pending := SELECT_ENTRY_Pending(bson.M{}, bson.M{})

	for _, entry := range pending {
		DELETE_ENTRY_Pending(entry)
	}
}

func HTTP_JSON_IdentityConfidence_Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	results := []Entity{}

	results = append(results, SELECT_ENTRY_Inventory(bson.M{"_id": system.EncodeID(id)}, bson.M{})...)
	results = append(results, SELECT_ENTRY_Pending(bson.M{"_id": system.EncodeID(id)}, bson.M{})...)

	if len(results) != 1 {
		system.Warning("Incorrect number of results returned for IdentityConfidence IP", true)
		w.Write([]byte("404/Failure"))
		return
	}

	json.NewEncoder(w).Encode(CalcIdentityConfidenceScore(results[0]))
}

/* HTTP_JSON_GetByID returns a CMDB element based on the {id} */
func HTTP_JSON_GetDateTimeGraph(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	entry := SELECT_Entities_Joined(bson.M{"_id": system.EncodeID(id)}, bson.M{})

	if len(entry) == 1 {
		json.NewEncoder(w).Encode(CalcTimeGraph(entry[0]))
	}
}

/*
	HTTP_JSON_Restore restores the system settings and databases
	to factory defaults.
*/
func HTTP_JSON_Restore(w http.ResponseWriter, r *http.Request) {
	system.Log("Restoring CMDB to factory defaults", true)
	system.CMDB_Inventory_DB.Drop(context.Background()) // Drop inventory
	system.CMDB_Pending_DB.Drop(context.Background())   // Drop pending

	Init() // Restore capabilities

	w.Write([]byte("200/Done"))
}

func HTTP_None_ResolveLocalIdentityProposition(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("ID")
	value := r.PostFormValue("Value")

	valueIdx, err := strconv.Atoi(value)
	system.Error("Couldn't convert proposition value ID to int", err)

	for i, prop := range proposition.Propositions {
		if prop.ID == id {
			setLocalIdentityEntry(proposition.Propositions[i].Predicates[valueIdx].Value)
			proposition.Pop(valueIdx)
			break
		}
	}
}

func HTTP_None_ResolveIPConflict(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("ID")
	value := r.PostFormValue("Value")

	valueIdx, err := strconv.Atoi(value)
	system.Error("Couldn't convert proposition value ID to int", err)

	for _, prop := range proposition.Propositions {
		if prop.ID == id {
			resolveIPConflict(ConflictActions(prop.Predicates[valueIdx].Value), prop.GetEvidenceValue("Conflict IP"))
			proposition.Pop(valueIdx)
			break
		}
	}
}
