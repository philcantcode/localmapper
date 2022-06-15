package cookbook

import (
	"time"

	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cookbook struct {
	ID         primitive.ObjectID `bson:"_id"`
	CCBI       string             // Common CookBook Identifier
	Label      string
	Desc       string
	CCIs       []string   // List of Common Capability IDs to run
	SearchKeys []string   // List of Labels {Port, IP, OS} to search for
	Schedule   []Schedule // Schedule info
	Hidden     bool       // Hidden from the user
}

/*
	FirstTimeSetup sets up any defaults necessary to run.
	All items should perform checks so they don't corrupt
	the database on subsiquent runs. Only when the database
	is empty should the initial setups run.
*/
func FirstTimeSetup() {
	localHostID := Cookbook{
		CCBI:  "ccbi:discovery:local-hosts:default", // Common CookBook Identifier: ccbi:<category>:<label>:<edition>
		Label: "Local Host Identification",
		Desc:  "Gather Hostname, Ports, OS, MAC, etc, about a local host.",
		CCIs: []string{
			"cci:nmap:os-ident:default",
			"cci:nmap:sys-dns:default",
		},
		SearchKeys: []string{"Ports", "OS", "OSGen", "OSVendor", "MAC", "HostName"},
		Hidden:     false,
		Schedule: []Schedule{
			{
				Label:    "Inventory Discovery",
				Desc:     "Once an item is in the inventory, scan the devices for info reguleraly.",
				Duration: time.Minute * 30,
				TargetGroups: []cmdb.CMDBType{
					cmdb.ENDPOINT,
					cmdb.SERVER,
					cmdb.PENDING,
				},
			},
		},
	}

	pingNonDefaultVLANs := Cookbook{
		CCBI:  "ccbi:discovery:ping-sweep-vlans:exclude-private-ranges", // Common CookBook Identifier: ccbi:<category>:<label>:<edition>
		Label: "Ping all VLANs",
		Desc:  "Pings all known VLAN excluding the three big private ranges.",
		CCIs: []string{
			"cci:nmap:ping-sweep:default",
		},
		SearchKeys: []string{},
		Hidden:     true,
		Schedule: []Schedule{
			{
				Label:    "Ping Sweep",
				Desc:     "Ping all known VLANs, excluding the three big private ranges.",
				Duration: time.Minute * 30,
				TargetGroups: []cmdb.CMDBType{
					cmdb.VLAN,
				},
				ExclusionList: []Exclusion{
					{
						Value:    "10.0.0.0/8",
						DataType: system.CIDR,
					},
					{
						Value:    "172.16.0.0/12",
						DataType: system.CIDR,
					},
					{
						Value:    "192.168.0.0/16",
						DataType: system.CIDR,
					},
				},
			},
		},
	}

	// pingSelf := Cookbook{
	// 	CCBI:  "ccbi:test:ping-self:default", // Common CookBook Identifier: ccbi:<category>:<label>:<edition>
	// 	Label: "Test Ping Self",
	// 	Desc:  "Test Cookbook.",
	// 	CCIs: []string{
	// 		"cci:nmap:ping-sweep:default",
	// 	},
	// 	SearchKeys: []string{},
	// 	Hidden:     true,
	// 	Schedule: []Schedule{
	// 		{
	// 			Label:    "Test Ping Self",
	// 			Desc:     "Test ping self every few seconds",
	// 			Duration: time.Second * 5,
	// 			TargetDevices: []string{
	// 				"62a71cce4880e441193254a4",
	// 			},
	// 		},
	// 	},
	// }

	if len(SELECT_Cookbook(bson.M{}, bson.M{})) == 0 {
		//INSERT_Cookbook(pingSelf)
		INSERT_Cookbook(localHostID)
		INSERT_Cookbook(pingNonDefaultVLANs)
	}
}
