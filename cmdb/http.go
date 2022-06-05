package cmdb

import (
	"encoding/json"
	"net/http"
)

func HTTP_JSON_GetSelf(w http.ResponseWriter, r *http.Request) {

	cmdbs := SELECT_CMDBItem_All()

	for _, cmdb := range cmdbs {
		tag, exists := cmdb.InfoTags["identity"]

		if exists && tag == "local" {
			json.NewEncoder(w).Encode(cmdb)
			return
		}
	}

	json.NewEncoder(w).Encode(CMDBItem{})
}

func HTTP_JSON_GetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SELECT_CMDBItem_All())
}
