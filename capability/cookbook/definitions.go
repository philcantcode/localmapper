package cookbook

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cookbook struct {
	ID         primitive.ObjectID `bson:"_id"`
	CCBI       string             // Common CookBook Identifier
	Label      string
	Desc       string
	CCIs       []string // List of Common Capability IDs to run
	SearchKeys []string // List of Labels {Port, IP, OS} to search for
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
	}

	if len(SELECT_Cookbook(bson.M{}, bson.M{})) == 0 {
		INSERT_Cookbook(localHostID)
	}
}
