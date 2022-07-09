package cmdb

import (
	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
0 = VLAN
1 = SERVER
2 = ENDPOINT
3 = UNKNOWN
*/
type CMDBType int

const (
	VLAN CMDBType = iota
	SERVER
	ENDPOINT
	PENDING
	SELF
)

/* OSI Layer:
   7 - Firewalls, IDS
   6 -
   5 -
   4 - Firewalls (some)
   3 - Routers, L3 switches
   2 - L2 switches, Bridges
   1 - Hubs, repeaters, modems */
type Entity struct {
	ID          primitive.ObjectID `bson:"_id"`
	Label       string
	Description string
	CMDBType    CMDBType
	OSILayer    int         // 1 - 7
	DateSeen    []string    //[] Array of dates seen
	UsrTags     []EntityTag // [Project-X, Bob's Server] etc
	SysTags     []EntityTag // [IP: xxx, MAC: xxx, URL: xxx] etc
	//TODO: implement histories for tracking changes
}

type EntityTag struct {
	Label       string
	Description string
	DataType    system.DataType
	Values      []string
}

type TimeGraph struct {
	Keys   []string
	Values []int
}

/*
	PushToFront moves a particular tag to the front
	of the list
*/
func (tag EntityTag) PushToFront(value string) EntityTag {

	newValues := []string{}

	for idx, val := range tag.Values {
		if val == value {
			newValues = append(newValues, val)

			for i, v := range tag.Values {
				if i == idx {
					continue
				}

				newValues = append(newValues, v)
			}

			tag.Values = newValues
			return tag
		}
	}

	system.Warning("Couldn't find tag value to PushToFront", false)

	return tag
}

type ConflictActions string

const (
	Action_MERGE_INTO_INVENTORY   ConflictActions = "Merge Into Inventory"
	Action_MERGE_INTO_PENDING     ConflictActions = "Merge Into Pending"
	Action_DELETE_PENDING_ENTRY   ConflictActions = "Delete Pending Entry"
	Action_DELETE_INVENTORY_ENTRY ConflictActions = "Delete Inventory Entry"
)
