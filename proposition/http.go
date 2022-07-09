package proposition

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

/* HTTP_JSON_GetDefaultGatewayIP both the deafult IP and the Gateway */
func HTTP_JSON_GetPropositions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Propositions)
}

func HTTP_None_Refresh(w http.ResponseWriter, r *http.Request) {
	// Init()
	fmt.Println("INITIALISE REFRESH CODE HERE")
}

func HTTP_INT_GetPropositionCount(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(strconv.Itoa(len(Propositions))))
}
