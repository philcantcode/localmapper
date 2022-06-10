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
)

/* OSI Layer:
   7 - Firewalls, IDS
   6 -
   5 -
   4 - Firewalls (some)
   3 - Routers, L3 switches
   2 - L2 switches, Bridges
   1 - Hubs, repeaters, modems */
type Entry struct {
	ID       primitive.ObjectID `bson:"_id"`
	Label    string
	Desc     string
	CMDBType CMDBType
	OSILayer int        // 1 - 7
	DateSeen []string   //[] Array of dates seen
	UsrTags  []EntryTag // [Project-X, Bob's Server] etc
	SysTags  []EntryTag // [IP: xxx, MAC: xxx, URL: xxx] etc
	//TODO: implement histories for tracking changes
}

type EntryTag struct {
	Label    string
	Desc     string
	DataType system.DataType
	Values   []string
}

type TimeGraph struct {
	Keys   []string
	Values []int
}
