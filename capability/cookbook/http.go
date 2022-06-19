package cookbook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
)

func HTTP_JSON_Run_Cookbook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ccbi := params["ccbi"]
	id := params["id"]

	cookbooks := SELECT_Cookbook(bson.M{"ccbi": ccbi}, bson.M{})
	entries := cmdb.SELECT_ENTRY_Joined(bson.M{"_id": system.EncodeID(id)}, bson.M{})

	if len(cookbooks) != 1 {
		system.Force("Wrong number of cookbooks selected", true)
		return
	}

	if len(entries) != 1 {
		system.Force(fmt.Sprintf("Wrong number of entries selected: %d", len(entries)), true)
		return
	}

	ExecuteCookbook(cookbooks[0], entries[0].ID)

	w.Write([]byte("200/Done"))
}

func HTTP_JSON_GetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SELECT_Cookbook(bson.M{}, bson.M{}))
}

/*
	HTTP_JSON_Restore restores the system settings and databases
	to factory defaults.
*/
func HTTP_JSON_Restore(w http.ResponseWriter, r *http.Request) {
	system.Log("Restoring cookbooks to factory defaults", true)

	DELETE_Cookbook(bson.M{})
	system.Core_Cookbooks_DB.Drop(context.Background()) // Drop cookbooks

	FirstTimeSetup() // Restore cookbooks

	w.Write([]byte("200/Done"))
}

func HTTP_JSON_GetNew(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Cookbook{CCIs: []string{}, SearchKeys: []string{}, Schedule: []Schedule{}})
}
