package webhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/cookbook"
	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
)

type CookbooksHandler struct {
}

var Cookbooks = CookbooksHandler{}

func (cb *CookbooksHandler) HTTP_JSON_Run_Cookbook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ccbi := params["ccbi"]
	id := params["id"]

	cookbooks := cookbook.SELECT_Cookbook(bson.M{"ccbi": ccbi}, bson.M{})
	entries := cmdb.SELECT_Entities_Joined(bson.M{"_id": system.EncodeID(id)}, bson.M{})

	if len(cookbooks) != 1 {
		system.Warning("Wrong number of cookbooks selected", true)
		return
	}

	if len(entries) != 1 {
		system.Warning(fmt.Sprintf("Wrong number of entries selected: %d", len(entries)), true)
		return
	}

	cookbooks[0].RunBookOnEntity(entries[0].ID)

	w.Write([]byte("200/Done"))
}

func (cb *CookbooksHandler) HTTP_JSON_GetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(cookbook.SELECT_Cookbook(bson.M{}, bson.M{}))
}

/*
	HTTP_JSON_Restore restores the system settings and databases
	to factory defaults.
*/
func (cb *CookbooksHandler) HTTP_JSON_Restore(w http.ResponseWriter, r *http.Request) {
	cookbook.Restore()
	w.Write([]byte("200/Done"))
}

func (cb *CookbooksHandler) HTTP_JSON_GetNew(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(cookbook.Cookbook{CCIs: []string{}, SearchKeys: []string{}, Schedule: []cookbook.Schedule{}})
}
