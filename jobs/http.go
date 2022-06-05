package jobs

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func HTTP_JSON_GetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(FILTER_JobSpec(bson.M{}, bson.M{}))
}
