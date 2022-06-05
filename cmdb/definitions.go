package cmdb

import "go.mongodb.org/mongo-driver/bson/primitive"

/* OSI Layer:
   7 - Firewalls, IDS
   6 -
   5 -
   4 - Firewalls (some)
   3 - Routers, L3 switches
   2 - L2 switches, Bridges
   1 - Hubs, repeaters, modems */
type CMDBItem struct {
	ID         primitive.ObjectID `bson:"_id"`
	Label      string
	Desc       string
	OSILayer   int               // 1 - 7
	DateSeen   []string          //[] Array of dates seen
	StatusTags map[string]string // [Stopped, Running] etc
	UserTags   map[string]string // [Project-X, Bob's Server] etc
	InfoTags   map[string]string // [IP: xxx, MAC: xxx, URL: xxx] etc
}

type Vlan struct {
	ID          int `bson:"_id"`
	Name        string
	Description string
	HighIP      string
	LowIP       string
	Tags        string
}
