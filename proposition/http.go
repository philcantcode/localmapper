package proposition

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/philcantcode/localmapper/system"
)

/* HTTP_None_Resolve runs when the user Accepts the default in the web gui
   this returns an ID*/
func HTTP_None_Resolve(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("ID")
	value := r.PostFormValue("Value")

	valueIdx, err := strconv.Atoi(value)
	system.Error("Couldn't convert proposition value ID to int", err)

	for i, prop := range propositions {
		if prop.ID == id {
			propositions[i].Predicates[valueIdx].Chosen = true
			propositions[i].resolve()

			break
		}
	}
}

/* HTTP_JSON_GetDefaultGatewayIP both the deafult IP and the Gateway */
func HTTP_JSON_GetPropositions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(propositions)
}

func HTTP_None_Refresh(w http.ResponseWriter, r *http.Request) {
	Init()
}
