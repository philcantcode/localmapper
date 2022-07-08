package proposition

import (
	"encoding/json"
	"net/http"
)

/* HTTP_None_Process runs when the user Accepts the default in the web gui
   this returns an ID*/
func HTTP_None_Process(w http.ResponseWriter, r *http.Request) {
	propStr := r.PostFormValue("Proposition")

	// Set the result to accept by default
	var proposition = Proposition{}
	json.Unmarshal([]byte(propStr), &proposition)

	proposition.process()
}

/* HTTP_JSON_GetDefaultGatewayIP both the deafult IP and the Gateway */
func HTTP_JSON_GetPropositions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(propositions)
}
