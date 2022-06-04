package propositions

import (
	"encoding/json"
	"net/http"

	"github.com/philcantcode/localmapper/database"
)

/* HTTP_JSON_GetDefaultGatewayIP both the deafult IP and the Gateway */
func HTTP_JSON_GetPropositions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(database.SelectAllPropositions())
}
