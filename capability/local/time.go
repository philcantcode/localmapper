package local

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/philcantcode/localmapper/utils"
)

type DateTime struct {
	DDMMYYYY string
	HHMMSS   string
	DateTime string
}

func GetDateTime() DateTime {
	dt := time.Now()
	dts := DateTime{}

	dts.DDMMYYYY = dt.Format("02-01-2006")
	dts.HHMMSS = dt.Format("15:04:05")
	dts.DateTime = dt.Format("02-01-2006 15:04:05")

	utils.Log("Returning the date & time.", false)

	return dts
}

/* HTTP_JSON_GetDateTime returns datetime */
func HTTP_JSON_GetDateTime(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(GetDateTime())
}
