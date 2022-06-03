package network

import (
	"github.com/philcantcode/localmapper/capabilities/nmap"
	"github.com/philcantcode/localmapper/database"
	"go.mongodb.org/mongo-driver/bson"
)

// ListAllAddresses finds all unique IP addresses from the database
// and returns a list
func ListAllAddresses() []nmap.Address {
	var returnArray []nmap.Address

	results := database.FilterNetworkNmap(
		bson.M{},
		bson.M{
			"hosts": bson.M{
				"addresses": 1,
			},
		})

	for _, collection := range results {
		for _, host := range collection.Hosts {
			for _, address := range host.Addresses {
				add := true
				for _, logged := range returnArray {
					if logged.Addr == address.Addr {
						add = false
						break
					}
				}

				if add {
					returnArray = append(returnArray, address)
				}
			}
		}
	}

	return returnArray
}
