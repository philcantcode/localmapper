package proposition

import (
	"github.com/philcantcode/localmapper/capability/local"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/utils"
)

func processProposition(proposition Proposition) {
	switch proposition.Type {
	case "local-net-iface":
		sysTags := []cmdb.EntityTag{}
		usrTags := []cmdb.EntityTag{}

		sysTags = append(sysTags, cmdb.EntityTag{Label: "Verified", DataType: system.DataType_BOOL, Values: []string{"1"}})
		sysTags = append(sysTags, cmdb.EntityTag{Label: "Identity", DataType: system.DataType_STRING, Values: []string{"local"}})

		sysTags = append(sysTags, cmdb.EntityTag{Label: "IP", DataType: system.DataType_IP, Values: []string{proposition.Predicate.Value}})

		for _, net := range local.GetNetworkAdapters() {
			if net.IP == proposition.Predicate.Value {
				if net.MAC != "" {
					sysTags = append(sysTags, cmdb.EntityTag{Label: "MAC", DataType: system.DataType_MAC, Values: []string{net.MAC}})
				}

				if net.MAC6 != "" {
					sysTags = append(sysTags, cmdb.EntityTag{Label: "MAC6", DataType: system.DataType_MAC6, Values: []string{net.MAC6}})
				}

				if net.Label != "" {
					sysTags = append(sysTags, cmdb.EntityTag{Label: "NetAdapter", DataType: system.DataType_STRING, Values: []string{net.Label}})
				}

				if net.IP6 != "" {
					sysTags = append(sysTags, cmdb.EntityTag{Label: "IP6", DataType: system.DataType_IP6, Values: []string{net.IP6}})
				}
			}
		}

		time := []string{utils.GetDateTime().DateTime}

		serverCMDB := cmdb.Entity{
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
