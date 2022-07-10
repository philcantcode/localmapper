package webhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
)

type SystemHandler struct {
}

var System = SystemHandler{}

func (sys *SystemHandler) HTTP_JSON_ExecuteAction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idstr := params["id"]

	id, err := strconv.Atoi(idstr)
	system.Fatal("Couldn't convert ID to int32: "+idstr, err)

	fmt.Println(id)
	fmt.Println("FIX ME::: " + idstr)
	//ExecuteCommand(id)

	w.Write([]byte("200/Done"))
}

/* HTTP_JSON_GetLogs returns all logs */
func (sys *SystemHandler) HTTP_JSON_GetLogs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(system.SELECT_LogEntry(bson.M{}, bson.M{}))
}

/*
	HTTP_JSON_Restore restores the system settings and databases
	to factory defaults.
*/
func (sys *SystemHandler) HTTP_JSON_Restore(w http.ResponseWriter, r *http.Request) {
	system.Restore()

	w.Write([]byte("200/Done"))
}
