package cmdb

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/capability/local"
	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func HTTP_JSON_GetLocal(w http.ResponseWriter, r *http.Request) {

	entries := SELECT_ENTRY_Inventory(bson.M{}, bson.M{})

	for _, entry := range entries {
		sysTag, exists, _ := FindSysTag("Identity", entry)
		if exists && utils.ArrayContains("local", sysTag.Values) {
			json.NewEncoder(w).Encode(entry)
			return
		}
	}

	json.NewEncoder(w).Encode(Entry{})
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
	utils.ErrorLog("Couldn't convert CMDBType to int32: "+typeStr, err, true)

	json.NewEncoder(w).Encode(SELECT_ENTRY_Inventory(bson.M{"cmdbtype": int32(typeNum)}, bson.M{}))
}

/* HTTP_JSON_GetByID returns a CMDB element based on the {id} */
func HTTP_JSON_GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	utils.Log("Inventory request made for: "+id, false)

	if id == "" {
		json.NewEncoder(w).Encode(SELECT_ENTRY_Inventory(bson.M{}, bson.M{}))
		return
	}

	// The inventory ID was a device
	device := SELECT_ENTRY_Inventory(bson.M{"_id": database.EncodeID(id)}, bson.M{})

	if len(device) > 0 {
		json.NewEncoder(w).Encode(device)
		return
	}

	json.NewEncoder(w).Encode(Entry{SysTags: []EntryTag{}, UsrTags: []EntryTag{}})
}

func HTTP_INSERT_Pending(w http.ResponseWriter, r *http.Request) {
	label := r.FormValue("Label")
	desc := r.FormValue("Desc")
	highIP := r.FormValue("HighIP")
	lowIP := r.FormValue("LowIP")
	cmdbType := r.FormValue("CMDBType")

	cmdbTypeInt, err := strconv.Atoi(cmdbType)
	utils.ErrorLog("Couldn't convert CMDBType to int", err, true)

	if !utils.ValidateIP(lowIP) {
		utils.ErrorContextLog("LowIP not valid in cmdb.HTTP_INSERT_Pending", true)
		return
	}

	if !utils.ValidateIP(highIP) {
		utils.ErrorContextLog("HighIP not valid in cmdb.HTTP_INSERT_Pending", true)
		return
	}

	if !utils.ValidateString(label) {
		utils.ErrorContextLog("Label not valid in cmdb.HTTP_INSERT_Pending", true)
		return
	}

	if !utils.ValidateString(desc) {
		utils.ErrorContextLog("Desc not valid in cmdb.HTTP_INSERT_Pending", true)
		return
	}

	highIpTag := EntryTag{Label: "LowIP", DataType: utils.IP, Values: []string{lowIP}}
	lowIpTag := EntryTag{Label: "HighIP", DataType: utils.IP, Values: []string{highIP}}
	entry := Entry{Label: label, Desc: desc, OSILayer: 2, CMDBType: CMDBType(cmdbTypeInt), DateSeen: []string{local.GetDateTime().DateTime}, SysTags: []EntryTag{lowIpTag, highIpTag}, UsrTags: []EntryTag{}}

	INSERT_ENTRY_Pending(entry)

	w.Write([]byte("200/Success"))
}

func HTTP_Pending_Approve(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("ID")

	pending := SELECT_ENTRY_Pending(bson.M{"_id": database.EncodeID(id)}, bson.M{})[0]
	pending.SysTags = append(pending.SysTags, EntryTag{Label: "Verified", DataType: utils.BOOL, Values: []string{"1"}})

	INSERT_ENTRY_Inventory(pending)
	DELETE_ENTRY_Pending(pending)
}

func HTTP_Pending_Deny(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("ID")

	pending := SELECT_ENTRY_Pending(bson.M{"_id": database.EncodeID(id)}, bson.M{})[0]

	DELETE_ENTRY_Pending(pending)
}

func HTTP_Pending_DenyAll(w http.ResponseWriter, r *http.Request) {
	pending := SELECT_ENTRY_Pending(bson.M{}, bson.M{})

	for _, entry := range pending {
		DELETE_ENTRY_Pending(entry)
	}
}
