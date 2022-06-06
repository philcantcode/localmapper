package proposition

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/philcantcode/localmapper/database"
	"go.mongodb.org/mongo-driver/bson"
)

/* HTTP_None_AcceptDefault runs when the user Accepts the default in the web gui
   this returns an ID*/
func HTTP_None_AcceptDefault(w http.ResponseWriter, r *http.Request) {
	ID := r.PostFormValue("ID")

	fmt.Println(ID)
	// Set the result to accept by default
	prop := SELECT_Propositions(bson.M{"_id": database.EncodeID(ID)}, bson.M{})[0]
	prop.Status = 1

	processProposition(prop)
}

/* HTTP_JSON_GetDefaultGatewayIP both the deafult IP and the Gateway */
func HTTP_JSON_GetPropositions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SELECT_Propositions(bson.M{}, bson.M{}))
}
