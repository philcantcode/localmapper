package gateway

import (
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/tools/searchsploit"
	"go.mongodb.org/mongo-driver/bson"
)

func FindEntityVulnerabilities(entityID string) []searchsploit.ExploitDB {
	entity := cmdb.SELECT_Entities_Joined(bson.M{"_id": system.EncodeID(entityID)}, bson.M{})
	vulnreabilities := []searchsploit.ExploitDB{}

	if len(entity) != 1 {
		system.Warning("Incorrect number of entities returned in find vulns", true)
	}

	// Join all tags
	allTags := entity[0].SysTags
	allTags = append(allTags, entity[0].UsrTags...)

	// Loop over SysTags and find vulnerabilities
	for _, tag := range allTags {
		for _, val := range tag.Values {
			// Search for vulnerabilities
			vulns := searchsploit.Select(bson.M{"search": val}, bson.M{})

			// Ensure they aren't already in the list
			for _, vuln := range vulns {
				exists := false

				for _, existingVulns := range vulnreabilities {
					if vuln.Search == existingVulns.Search {
						exists = true
					}
				}

				if !exists {
					vulnreabilities = append(vulnreabilities, vuln)
				}
			}
		}
	}

	return vulnreabilities
}
