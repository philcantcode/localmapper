package gateway

import (
	"github.com/philcantcode/localmapper/capability"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
)

/*
	Given a CMDB Entity ID, finds a set of compatible
	Capabilities
*/
func FindCompatibleCapabilitiesForEntity(cmbdID string) []capability.Capability {
	result := []capability.Capability{}
	entries := cmdb.SELECT_Entities_Joined(bson.M{"_id": system.EncodeID(cmbdID)}, bson.M{})

	if len(entries) != 1 {
		system.Warning("Too many results returned in FindCompatibleCapabilitiesForEntity by ID", true)
	}

	caps := capability.SELECT_Capability(bson.M{}, bson.M{})

	for _, cap := range caps {
		isCompatible, parsedCap := cap.CheckCompatability(entries[0])

		if isCompatible {
			result = append(result, parsedCap)
		}
	}

	return result
}

/*
	Given a Capability ID, a CMDB Entity ID and a set of
	KEY(string):VALUE(int) option pairs - which represent
	a Flag and OptionIDX (?) attempts to run the capability
*/
func LaunchCapabilityForEntity(capID string, entityID string, options map[string]int) {
	cap := capability.SELECT_Capability(bson.M{"_id": system.EncodeID(capID)}, bson.M{})[0]
	entries := cmdb.SELECT_Entities_Joined(bson.M{"_id": system.EncodeID(entityID)}, bson.M{})

	if len(entries) != 1 {
		system.Warning("Too many results returned in HTTP_JSON_GetCMDBCompatible", true)
		return
	}

	for key, val := range options {
		for idx, param := range cap.Command.Params {
			if param.Flag == key {
				cap.Command.Params[idx].Value = param.Options[val].Value
			}
		}
	}

	isCompatible, parsedCap := cap.CheckCompatability(entries[0])

	if !isCompatible {
		return
	}

	parsedCap.LaunchLifecycleManager()
}
