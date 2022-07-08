package proposition

import (
	"time"

	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/local"
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func setupLocalIPIdentity() {
	entries := cmdb.SELECT_ENTRY_Inventory(bson.M{"systags.label": "Identity", "systags.values": "local"}, bson.M{})

	if len(entries) > 1 {
		system.Warning("Too many propositions for self idetntity returned", true)
	}

	// Self identity already setup
	if len(entries) == 1 {
		return
	}

	// No identity setup
	// Create the new proposition to guess the local IP

	proposition := Proposition{
		Type:        Proposition_Local_Identity,
		DateTime:    time.Now(),
		Description: "Please choose the IP address for this server.",
		Predicates: []Predicate{
			{
				Label:    "Default IP",
				Value:    local.GetDefaultIPGateway().DefaultIP,
				DataType: system.DataType_IP,
			},
		},
	}

	for _, ip := range local.GetNetworkAdapters() {
		proposition.Predicates = append(proposition.Predicates, Predicate{
			Label:    "Adapter IP",
			Value:    ip.IP,
			DataType: system.DataType_IP,
		})
	}

	proposition.ID = utils.HashStruct(proposition)
	propositions = append(propositions, proposition)
}

func (proposition Proposition) resolveLocalIPIdentity() {
	sysTags := []cmdb.EntityTag{}
	usrTags := []cmdb.EntityTag{}
	chosen := proposition.getChosen()

	sysTags = append(sysTags, cmdb.EntityTag{Label: "Verified", DataType: system.DataType_BOOL, Values: []string{"1"}})
	sysTags = append(sysTags, cmdb.EntityTag{Label: "Identity", DataType: system.DataType_STRING, Values: []string{"local"}})
	sysTags = append(sysTags, cmdb.EntityTag{Label: "IP", DataType: system.DataType_IP, Values: []string{chosen.Value}})

	for _, net := range local.GetNetworkAdapters() {
		if net.IP == chosen.Value {
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
		Label:       "Local-Mapper Server",
		OSILayer:    7,
		Description: "The local-mapper backend server.",
		DateSeen:    time,
		CMDBType:    cmdb.SERVER,
		UsrTags:     usrTags,
		SysTags:     sysTags,
	}

	serverCMDB.InsertInventory()
}
