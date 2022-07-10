package webhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/philcantcode/localmapper/proposition"
)

type PropositionHandler struct {
}

var Proposition = PropositionHandler{}

/* HTTP_JSON_GetDefaultGatewayIP both the deafult IP and the Gateway */
func (prop *PropositionHandler) HTTP_JSON_GetPropositions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(proposition.Proposition{})
}

func (prop *PropositionHandler) HTTP_None_Refresh(w http.ResponseWriter, r *http.Request) {
	// Init()
	fmt.Println("INITIALISE REFRESH CODE HERE")
}

func (prop *PropositionHandler) HTTP_INT_GetPropositionCount(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(strconv.Itoa(len(proposition.Propositions))))
}
