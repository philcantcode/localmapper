package local

import (
	"encoding/json"
	"net/http"
	"time"
)

type DateTime struct {
	DDMMYYYY string
	HHMMSS   string
}

func GetDateTime() DateTime {
	dt := time.Now()
	dts := DateTime{}

	dts.DDMMYYYY = dt.Format("02-01-2006")
	dts.HHMMSS = dt.Format("15:04:05")

	return dts
}

/* HTTP_JSON_GetDateTime returns datetime */
func HTTP_JSON_GetDateTime(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(GetDateTime())
}
