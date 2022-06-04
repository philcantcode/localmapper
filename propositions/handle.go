package propositions

import (
	"net/http"
	"strconv"

	"github.com/philcantcode/localmapper/capabilities/local"
	"github.com/philcantcode/localmapper/core"
	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
)

/* ProcessAcceptDefaults runs when the user Accepts the default in the web gui
   this returns an ID*/
func ProcessAcceptDefaults(w http.ResponseWriter, r *http.Request) {
	ID := r.PostFormValue("ID")
	idInt, err := strconv.Atoi(ID)
	utils.ErrorLog("Couldn't convert ID to integer while running ProcessAcceptDefaults", err, true)

	// Set the result to accept by default
	prop := database.SelectPropositionByID(idInt)
	prop.Status = 0
	prop.Correction = prop.Proposition

	processProposition(prop)
}

func processProposition(proposition core.Proposition) {
	switch proposition.Type {
	case "local-net-iface":
		statusTags := make(map[string]string)
		infoTags := make(map[string]string)

		statusTags["Verified"] = "1"
		infoTags["IP"] = proposition.Correction.Value
		infoTags["identity"] = "local"

		for _, net := range local.GetNetworkAdapters() {
			if net.IP == proposition.Correction.Value {
				infoTags["MAC"] = net.MAC
				infoTags["MAC6"] = net.MAC6
				infoTags["NetAdapter"] = net.Name
				infoTags["IP6"] = net.IP6
			}
		}

		time := []string{local.GetDateTime().DateTime}

		serverCMDB := core.CMDBItem{OSILayer: 7, Description: "local-mapper Server", StatusTags: statusTags, InfoTags: infoTags, DateSeen: time}
		database.InsertCMDBItem(serverCMDB)

		database.SetPropositionStatusByID(proposition.ID, 1)
	}
}
