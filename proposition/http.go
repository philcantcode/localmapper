package proposition

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/philcantcode/localmapper/utils"
)

/* ProcessAcceptDefaults runs when the user Accepts the default in the web gui
   this returns an ID*/
func ProcessAcceptDefaults(w http.ResponseWriter, r *http.Request) {
	ID := r.PostFormValue("ID")
	idInt, err := strconv.Atoi(ID)
	utils.ErrorLog("Couldn't convert ID to integer while running ProcessAcceptDefaults", err, true)

	// Set the result to accept by default
	prop := SelectPropositionByID(idInt)
	prop.Status = 0
	prop.Correction = prop.Proposition

	processProposition(prop)
}

/* HTTP_JSON_GetDefaultGatewayIP both the deafult IP and the Gateway */
func HTTP_JSON_GetPropositions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SelectAllPropositions())
}
