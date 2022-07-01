package nmap

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func HTTP_SELECT_Nmap(w http.ResponseWriter, r *http.Request) {
	filter := r.FormValue("filter")
	projection := r.FormValue("projection")

	results := SELECT(, bson.M{})

	json.NewEncoder(w).Encode(results)
}
