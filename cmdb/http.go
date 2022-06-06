package cmdb

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func HTTP_JSON_GetLocal(w http.ResponseWriter, r *http.Request) {

	entries := SELECT_ENTRY(bson.M{}, bson.M{})

	for _, entry := range entries {
		sysTag, exists := FindSysTag("identity", entry)
		if exists && utils.ArrayContains("local", sysTag.Values) {
			json.NewEncoder(w).Encode(entry)
			return
		}
	}

	json.NewEncoder(w).Encode(Entry{})
}

func HTTP_JSON_GetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SELECT_ENTRY(bson.M{}, bson.M{}))
}

func HTTP_JSON_GetByType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	typeStr := params["type"]
	typeNum, err := strconv.Atoi(typeStr)
	utils.ErrorLog("Couldn't convert CMDBType to int32: "+typeStr, err, true)

	json.NewEncoder(w).Encode(SELECT_ENTRY(bson.M{"cmdbtype": int32(typeNum)}, bson.M{}))
}

/* HTTP_JSON_GetByID returns a CMDB element based on the {id} */
func HTTP_JSON_GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	utils.Log("Inventory request made for: "+id, false)

	if id == "" {
		json.NewEncoder(w).Encode(SELECT_ENTRY(bson.M{}, bson.M{}))
		return
	}

	// The inventory ID was a device
	device := SELECT_ENTRY(bson.M{"_id": database.EncodeID(id)}, bson.M{})

	if len(device) > 0 {
		json.NewEncoder(w).Encode(device)
		return
	}

}
