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

		sysTags = append(sysTags, cmdb.EntryTag{Label: "Verified", DataType: utils.BOOL, Values: []string{"1"}})
		sysTags = append(sysTags, cmdb.EntryTag{Label: "Identity", DataType: utils.STRING, Values: []string{"local"}})

		sysTags = append(sysTags, cmdb.EntryTag{Label: "IP", DataType: utils.IP, Values: []string{proposition.Predicate.Value}})

		for _, net := range local.GetNetworkAdapters() {
			if net.IP == proposition.Predicate.Value {
				if net.MAC != "" {
					sysTags = append(sysTags, cmdb.EntryTag{Label: "MAC", DataType: utils.MAC, Values: []string{net.MAC}})
				}

				if net.MAC6 != "" {
					sysTags = append(sysTags, cmdb.EntryTag{Label: "MAC6", DataType: utils.MAC6, Values: []string{net.MAC6}})
				}

				if net.Name != "" {
					sysTags = append(sysTags, cmdb.EntryTag{Label: "NetAdapter", DataType: utils.STRING, Values: []string{net.Name}})
				}

				if net.IP6 != "" {
					sysTags = append(sysTags, cmdb.EntryTag{Label: "IP6", DataType: utils.IP6, Values: []string{net.IP6}})
				}
			}
		}

		time := []string{local.GetDateTime().DateTime}

		serverCMDB := cmdb.Entry{
			Label:    "Local-Mapper Server (local)",
			OSILayer: 7,
			Desc:     "The local-mapper backend server.",
			DateSeen: time,
			CMDBType: cmdb.SERVER,
			UsrTags:  usrTags,
			SysTags:  sysTags}

		cmdb.INSERT_ENTRY_Inventory(serverCMDB)

		UPDATE_Proposition(proposition)
	}
}
