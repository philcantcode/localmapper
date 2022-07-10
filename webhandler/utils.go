package webhandler

import (
	"encoding/json"
	"net/http"

	"github.com/philcantcode/localmapper/utils"
)

type UtilsHandler struct {
}

var Utils = UtilsHandler{}

/* HTTP_JSON_GetDateTime returns datetime */
func (utls *UtilsHandler) HTTP_JSON_GetDateTime(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(utils.GetDateTime())
}
