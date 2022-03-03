package network

import (
	"github.com/philcantcode/localmapper/adapters/blueprint"
	"github.com/philcantcode/localmapper/application/database"
	"go.mongodb.org/mongo-driver/bson"
)

// ListAllAddresses finds all unique IP addresses from the database
// and returns a list
func ListAllAddresses() []blueprint.Address {
	var returnArray []blueprint.Address

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
