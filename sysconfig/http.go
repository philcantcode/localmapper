package sysconfig

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/utils"
)

func HTTP_JSON_Settings(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idstr := params["id"]

	id, err := strconv.Atoi(idstr)
	utils.ErrorLog("Couldn't convert ID to int32: "+idstr, err, true)

	ExecuteCommand(id)

	w.Write([]byte("200/Done"))
}
