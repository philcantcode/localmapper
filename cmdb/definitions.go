package cmdb

/* OSI Layer:
   7 - Firewalls, IDS
   6 -
   5 -
   4 - Firewalls (some)
   3 - Routers, L3 switches
   2 - L2 switches, Bridges
   1 - Hubs, repeaters, modems */
type CMDBItem struct {
	ID          int
	OSILayer    int      // 1 - 7
	DateSeen    []string //[] Array of dates seen
	Description string
	StatusTags  map[string]string // [Stopped, Running] etc
	UserTags    map[string]string // [Project-X, Bob's Server] etc
	InfoTags    map[string]string // [IP: xxx, MAC: xxx, URL: xxx] etc
}

type Vlan struct {
	ID          int    `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	HighIP      string `json:"HighIP"`
	LowIP       string `json:"LowIP"`
	Tags        string `json:"Tags"`
}
