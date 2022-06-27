package utils

import (
	"encoding/json"
	"net/http"
)

/* HTTP_JSON_GetDateTime returns datetime */
func HTTP_JSON_GetDateTime(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(GetDateTime())
}
