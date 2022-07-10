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

/*
	Given an IP, functional finds additional info and then
	creates an inventory entry for the IP.
*/
func SetLocalIdentityEntry(ip string) {
	sysTags := []EntityTag{}
	usrTags := []EntityTag{}

	sysTags = append(sysTags, EntityTag{Label: "Verified", DataType: system.DataType_BOOL, Values: []string{"1"}})
	sysTags = append(sysTags, EntityTag{Label: "Identity", DataType: system.DataType_STRING, Values: []string{"local"}})
	sysTags = append(sysTags, EntityTag{Label: "IP", DataType: system.DataType_IP, Values: []string{ip}})

	for _, net := range local.GetNetworkAdapters() {
		if net.IP == ip {
			if net.MAC != "" {
				sysTags = append(sysTags, EntityTag{Label: "MAC", DataType: system.DataType_MAC, Values: []string{net.MAC}})
			}

			if net.MAC6 != "" {
				sysTags = append(sysTags, EntityTag{Label: "MAC6", DataType: system.DataType_MAC6, Values: []string{net.MAC6}})
			}

			if net.Label != "" {
				sysTags = append(sysTags, EntityTag{Label: "NetAdapter", DataType: system.DataType_STRING, Values: []string{net.Label}})
			}

			if net.IP6 != "" {
				sysTags = append(sysTags, EntityTag{Label: "IP6", DataType: system.DataType_IP6, Values: []string{net.IP6}})
			}
		}
	}

	time := []string{utils.GetDateTime().DateTime}

	serverCMDB := Entity{
		Label:       "Local-Mapper Server",
		OSILayer:    7,
		Description: "The local-mapper backend server.",
		DateSeen:    time,
		CMDBType:    SERVER,
		UsrTags:     usrTags,
		SysTags:     sysTags,
	}

	serverCMDB.InsertInventory()
}

func ResolveIPConflict(action ConflictActions, ip string) {

	pending := SELECT_ENTRY_Pending(bson.M{"systags.label": "IP", "systags.values": ip}, bson.M{})[0]
	inventory := SELECT_ENTRY_Inventory(bson.M{"systags.label": "IP", "systags.values": ip}, bson.M{})[0]

	if action == Action_MERGE_INTO_INVENTORY {
		// Parse SysTags and join them
		for _, newTag := range pending.SysTags {
			_, found, i := inventory.FindSysTag(newTag.Label)

			if found {
				inventory.SysTags[i].Values = joinTagGroups(newTag.Label, inventory.SysTags[i].Values, newTag.Values)
			} else {
				inventory.SysTags = append(inventory.SysTags, newTag)
			}
		}

		// Parse SysTags and join them
		for _, newTag := range pending.UsrTags {
			_, found, i := inventory.FindUsrTag(newTag.Label)

			if found {
				inventory.UsrTags[i].Values = joinTagGroups(newTag.Label, inventory.UsrTags[i].Values, newTag.Values)
			} else {
				inventory.UsrTags = append(inventory.UsrTags, newTag)
			}
		}

		pending.UPDATE_ENTRY_Inventory()
		pending.DELETE_ENTRY_Pending()
		system.Log("Merged into Inventory", true)
	}

	if action == Action_MERGE_INTO_PENDING {
		// Parse SysTags and join them
		for _, newTag := range inventory.SysTags {
			_, found, i := pending.FindSysTag(newTag.Label)

			if found {
				pending.SysTags[i].Values = joinTagGroups(newTag.Label, pending.SysTags[i].Values, newTag.Values)
			} else {
				pending.SysTags = append(pending.SysTags, newTag)
			}
		}

		// Parse SysTags and join them
		for _, newTag := range inventory.UsrTags {
			_, found, i := pending.FindUsrTag(newTag.Label)

			if found {
				pending.UsrTags[i].Values = joinTagGroups(newTag.Label, pending.UsrTags[i].Values, newTag.Values)
			} else {
				pending.UsrTags = append(pending.UsrTags, newTag)
			}
		}

		pending.UPDATE_ENTRY_Inventory()
		DELETE_ENTRY_Inventory(inventory)
		system.Log("Merged into Pending", true)
	}

	if action == Action_DELETE_INVENTORY_ENTRY {
		DELETE_ENTRY_Inventory(inventory)
		system.Log("Deletd Inventory", true)
	}

	if action == Action_DELETE_PENDING_ENTRY {
		pending.DELETE_ENTRY_Pending()
		system.Log("Deletd Pending", true)
	}
}
