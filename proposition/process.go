package proposition

import (
	"github.com/philcantcode/localmapper/capability/local"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/utils"
)

func processProposition(proposition Proposition) {
	switch proposition.Type {
	case "local-net-iface":
		sysTags := []cmdb.EntryTag{}
		usrTags := []cmdb.EntryTag{}

		sysTags = append(sysTags, cmdb.EntryTag{Label: "verified", DataType: utils.Bool, Values: []string{"1"}})
		sysTags = append(sysTags, cmdb.EntryTag{Label: "identity", DataType: utils.Bool, Values: []string{"local"}})

		usrTags = append(usrTags, cmdb.EntryTag{Label: "IP", DataType: utils.IP, Values: []string{proposition.Predicate.Value}})

		for _, net := range local.GetNetworkAdapters() {
			if net.IP == proposition.Predicate.Value {
				usrTags = append(usrTags, cmdb.EntryTag{Label: "MAC", DataType: utils.MAC, Values: []string{net.MAC}})
				usrTags = append(usrTags, cmdb.EntryTag{Label: "MAC6", DataType: utils.MAC6, Values: []string{net.MAC6}})
				usrTags = append(usrTags, cmdb.EntryTag{Label: "NetAdapter", DataType: utils.String, Values: []string{net.Name}})
				usrTags = append(usrTags, cmdb.EntryTag{Label: "IP6", DataType: utils.IP6, Values: []string{net.IP6}})
			}
		}

		time := []string{local.GetDateTime().DateTime}

		serverCMDB := cmdb.Entry{
			Label:    "Local-Mapper Server (local)",
			OSILayer: 7,
			Desc:     "The local-mapper backend server.",
			DateSeen: time,
			UsrTags:  usrTags,
			SysTags:  sysTags}

		cmdb.INSERT_ENTRY(serverCMDB)

		UPDATE_Proposition(proposition)
	}
}
