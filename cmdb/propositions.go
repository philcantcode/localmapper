package cmdb

import (
	"fmt"

	"github.com/philcantcode/localmapper/local"
	"github.com/philcantcode/localmapper/proposition"
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type EntityConflicts struct {
	EntityOne Entity
	EntityTwo Entity
}

/*
	Creates a new proposition to describe an IP conflict between
	a device in the pending and inventory database.
*/
func InitIPConflict(pending Entity, inventory Entity) {
	tag, _, _ := inventory.FindSysTag("IP")
	ip := tag.Values[0]

	prop := proposition.Proposition{
		Type:        proposition.Proposition_IP_Identity_Conflict,
		DateTime:    utils.GetDateTime().DateTime,
		Description: "Resolve the conflict between two entities with the same IP address in the pending and inventory databases",
		Predicates: []proposition.Predicate{
			{
				Label:    "Action",
				DataType: system.DataType_STRING,
				Value:    string(Action_MERGE_INTO_INVENTORY),
			},
			{
				Label:    "Action",
				DataType: system.DataType_STRING,
				Value:    string(Action_MERGE_INTO_PENDING),
			},
			{
				Label:    "Action",
				DataType: system.DataType_STRING,
				Value:    string(Action_DELETE_INVENTORY_ENTRY),
			},
			{
				Label:    "Action",
				DataType: system.DataType_STRING,
				Value:    string(Action_DELETE_PENDING_ENTRY),
			},
		},
		Evidence: []proposition.Evidence{
			{
				Label: "Conflict IP",
				Value: ip,
			},
			{
				Label: "Pending Entity",
				Value: fmt.Sprintf("%+v", pending),
			},
			{
				Label: "Inventory Entity",
				Value: fmt.Sprintf("%+v", inventory),
			},
		},
	}

	prop.Push()
}

func InitLocalIdentityProp() {
	// Setup the self identity if doesn't exist
	entries := SELECT_ENTRY_Inventory(bson.M{"systags.label": "Identity", "systags.values": "local"}, bson.M{})

	if len(entries) > 1 {
		system.Warning("Too many propositions for self idetntity returned", true)
	}

	// Already exists
	if len(entries) == 1 {
		return
	}

	// There should only be one of them
	if proposition.CheckPropTypeExists(proposition.Proposition_Local_Identity) {
		return
	}

	// Create the new proposition to guess the local IP
	prop := proposition.Proposition{
		Type:        proposition.Proposition_Local_Identity,
		DateTime:    utils.GetDateTime().DateTime,
		Description: "Please choose the IP address for this server.",
		Predicates: []proposition.Predicate{
			{
				Label:    "Default IP",
				Value:    local.GetDefaultIPGateway().DefaultIP,
				DataType: system.DataType_IP,
			},
		},
	}

	for _, ip := range local.GetNetworkAdapters() {
		prop.Predicates = append(prop.Predicates, proposition.Predicate{
			Label:    "Adapter IP",
			Value:    ip.IP,
			DataType: system.DataType_IP,
		})
	}

	prop.Push()
}
