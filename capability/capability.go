package capability

import (
	"encoding/json"
	"fmt"

	"github.com/philcantcode/localmapper/capability/nmap"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/utils"
)

func ProcessCapability(capability Capability) []byte {
	switch capability.Type {
	case "nmap":
		nmapRun := nmap.Execute(ParamsToArray(capability.Command.Params))
		nmap.INSERT_Nmap(nmapRun)
		utils.PrintLog(utils.PrettyPrintToStr(nmapRun))

		result, err := json.Marshal(nmapRun)
		utils.ErrorLog("Couldn't marshal nmaprun", err, true)

		return result
	default:
		utils.ErrorForceFatal(fmt.Sprintf("No capability type to run in Capability.ProcessCapability: %v", capability))
		return nil
	}
}

func CMP_Entry_Capability(capability Capability, entry cmdb.Entry) (bool, Capability) {
	var success bool

	for k, capParam := range capability.Command.Params {
		success, capability.Command.Params[k] = CMP_CapabilityParam_Entry(capParam, entry.SysTags)

		if !success {
			return false, capability
		}
	}

	return true, capability
}

/*
Determines if given a capability param {"Value": "","DataType": 1, "Default": ""}
Is there any SysTags that can fulfil the Values
*/
func CMP_CapabilityParam_Entry(capParam Param, entryTags []cmdb.EntryTag) (bool, Param) {
	// For each: {DataType.CMDB, DataType.IP}
	for _, pType := range capParam.DataType {
		// If the value is already set, move on
		if capParam.Value != "" {
			return true, capParam
		}

		// Skip empty tags that don't require input
		if pType == utils.EMPTY {
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
