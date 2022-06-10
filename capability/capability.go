package capability

import (
	"encoding/json"
	"fmt"

	"github.com/philcantcode/localmapper/capability/nmap"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
)

func ExecuteCapability(capability Capability) []byte {
	system.Log(fmt.Sprintf("Executing Capability: %s\n", capability.Name), true)

	switch capability.Type {
	case "nmap":
		nmapRun := nmap.Execute(ParamsToArray(capability.Command.Params))
		nmap.INSERT_Nmap(nmapRun)

		result, err := json.Marshal(nmapRun)
		system.Error("Couldn't marshal nmaprun", err)

		return result
	default:
		system.Force(fmt.Sprintf("No capability type to run in Capability.ProcessCapability: %v", capability), true)
		return nil
	}
}

/*
	MatchEntryToCapability determines if a given entry can run a given capability
*/
func MatchEntryToCapability(capability Capability, entry cmdb.Entry) (bool, Capability) {
	var success bool

	for k, capParam := range capability.Command.Params {
		success, capability.Command.Params[k] = MatchParamToTag(capParam, entry.SysTags)

		if !success {
			return false, capability
		}
	}

	return true, capability
}

/*
	MatchParamToTag Determines if given a capability param {"Value": "","DataType": 1, "Default": ""}
	Is there any SysTags that can fulfil the Values
*/
func MatchParamToTag(capParam Param, entryTags []cmdb.EntryTag) (bool, Param) {
	// For each: {DataType.CMDB, DataType.IP}
	for _, pType := range capParam.DataType {
		// If the value is already set, move on
		if capParam.Value != "" {
			return true, capParam
		}

		// Skip empty tags that don't require input
		if pType == system.EMPTY {
			return true, capParam
		}

		for _, eTag := range entryTags {
			// The DataType match
			if pType == eTag.DataType {
				capParam.Value = eTag.Values[len(eTag.Values)-1]
				return true, capParam
			}
		}
	}

	return false, capParam
}
