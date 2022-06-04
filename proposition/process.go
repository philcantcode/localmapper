package proposition

import (
	"github.com/philcantcode/localmapper/capability/local"
	"github.com/philcantcode/localmapper/cmdb"
)

func processProposition(proposition Proposition) {
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

		serverCMDB := cmdb.CMDBItem{OSILayer: 7, Description: "local-mapper Server", StatusTags: statusTags, InfoTags: infoTags, DateSeen: time}
		cmdb.InsertCMDBItem(serverCMDB)

		SetPropositionStatusByID(proposition.ID, 1)
	}
}
