package jobs

import (
	"encoding/json"
	"net/http"
)

func HTTP_JSON_GetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SELECT_JobSpec_All())
}
