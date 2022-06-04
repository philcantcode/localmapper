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
	prop := SELECT_Proposition_ByID(idInt)
	prop.Status = 0
	prop.Correction = prop.Proposition

	processProposition(prop)
}

/* HTTP_JSON_GetDefaultGatewayIP both the deafult IP and the Gateway */
func HTTP_JSON_GetPropositions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SELECT_Propositions_All())
}

// HTTP_JSON_Refresh is called when the user visits http://server.com/propositions to referesh them
func HTTP_JSON_Refresh(w http.ResponseWriter, r *http.Request) {
	utils.Log("Refreshing propositions (http req)", false)

	SetupJobs()
	json.NewEncoder(w).Encode("Done")
}
