package cookbook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func HTTP_JSON_Run_Cookbook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ccbi := params["ccbi"]
	id := params["id"]

	cookbooks := SELECT_Cookbook(bson.M{"ccbi": ccbi}, bson.M{})
	entries := cmdb.SELECT_ENTRY_Joined(bson.M{"_id": database.EncodeID(id)}, bson.M{})

	if len(cookbooks) != 1 {
		utils.ErrorContextLog("Wrong number of cookbooks selected", true)
		return
	}

	if len(entries) != 1 {
		utils.ErrorContextLog(fmt.Sprintf("Wrong number of entries selected: %d\n", len(entries)), true)
		return
	}

	ExecuteCookbook(cookbooks[0], entries[0].ID)

	w.Write([]byte("200/Done"))
}

func HTTP_JSON_GetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SELECT_Cookbook(bson.M{}, bson.M{}))
}
