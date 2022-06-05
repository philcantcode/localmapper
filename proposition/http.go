package proposition

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

/* HTTP_None_AcceptDefault runs when the user Accepts the default in the web gui
   this returns an ID*/
func HTTP_None_AcceptDefault(w http.ResponseWriter, r *http.Request) {
	ID := r.PostFormValue("ID")

	fmt.Println(ID)
	// Set the result to accept by default
	prop := SELECT_Propositions(bson.M{"_id": database.ID_TO_Obj(ID)}, bson.M{})[0]
	prop.Status = 1

	processProposition(prop)
}

/* HTTP_JSON_GetDefaultGatewayIP both the deafult IP and the Gateway */
func HTTP_JSON_GetPropositions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SELECT_Propositions(bson.M{}, bson.M{}))
}

// HTTP_None_Refresh is called when the user visits http://server.com/propositions to referesh them
func HTTP_None_Refresh(w http.ResponseWriter, r *http.Request) {
	utils.Log("Refreshing propositions (http req)", false)

	SetupJobs()
}
