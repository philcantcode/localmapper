package proposition

import (
	"github.com/philcantcode/localmapper/capability/local"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
)

func processProposition(proposition Proposition) {
	switch proposition.Type {
	case "local-net-iface":
		sysTags := []cmdb.EntryTag{}
		usrTags := []cmdb.EntryTag{}

		sysTags = append(sysTags, cmdb.EntryTag{Label: "Verified", DataType: system.BOOL, Values: []string{"1"}})
		sysTags = append(sysTags, cmdb.EntryTag{Label: "Identity", DataType: system.STRING, Values: []string{"local"}})

		sysTags = append(sysTags, cmdb.EntryTag{Label: "IP", DataType: system.IP, Values: []string{proposition.Predicate.Value}})

		for _, net := range local.GetNetworkAdapters() {
			if net.IP == proposition.Predicate.Value {
				if net.MAC != "" {
					sysTags = append(sysTags, cmdb.EntryTag{Label: "MAC", DataType: system.MAC, Values: []string{net.MAC}})
				}

				if net.MAC6 != "" {
					sysTags = append(sysTags, cmdb.EntryTag{Label: "MAC6", DataType: system.MAC6, Values: []string{net.MAC6}})
				}

				if net.Label != "" {
					sysTags = append(sysTags, cmdb.EntryTag{Label: "NetAdapter", DataType: system.STRING, Values: []string{net.Label}})
				}

				if net.IP6 != "" {
					sysTags = append(sysTags, cmdb.EntryTag{Label: "IP6", DataType: system.IP6, Values: []string{net.IP6}})
				}
			}
		}

		time := []string{local.GetDateTime().DateTime}

		serverCMDB := cmdb.Entry{
			Label:       "Local-Mapper Server (local)",
			OSILayer:    7,
			Description: "The local-mapper backend server.",
			DateSeen:    time,
			CMDBType:    cmdb.SERVER,
			UsrTags:     usrTags,
			SysTags:     sysTags}

		cmdb.UpdateOrInsert(serverCMDB)

		UPDATE_Proposition(proposition)
	}
}
