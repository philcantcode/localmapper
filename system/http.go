package system

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func HTTP_JSON_ExecuteAction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idstr := params["id"]

	id, err := strconv.Atoi(idstr)
	Fatal("Couldn't convert ID to int32: "+idstr, err)

	fmt.Println(id)
	fmt.Println("FIX ME::: " + idstr)
	//ExecuteCommand(id)

	w.Write([]byte("200/Done"))
}

/* HTTP_JSON_GetLogs returns all logs */
func HTTP_JSON_GetLogs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SELECT_LogEntry(bson.M{}, bson.M{}))
}

/*
	HTTP_JSON_Restore restores the system settings and databases
	to factory defaults.
*/
func HTTP_JSON_Restore(w http.ResponseWriter, r *http.Request) {
	Log("Restoring system settings to factory defaults", true)

	System_Logs_DB.Drop(context.Background()) // Drop the logs table
	DELETE_Settings_All()                     // Delete all settings

	Init() // Perform first time setup

	w.Write([]byte("200/Done"))
}
