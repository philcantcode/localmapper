package system

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func HTTP_JSON_Settings(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idstr := params["id"]

	id, err := strconv.Atoi(idstr)
	Fatal("Couldn't convert ID to int32: "+idstr, err)

	fmt.Println(id)
	fmt.Println("FIX ME::: " + idstr)
	//ExecuteCommand(id)

	w.Write([]byte("200/Done"))
}
