package system

import (
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
