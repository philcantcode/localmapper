package capability

import (
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/utils"
)

/*
	CheckCompatability takes a given entry and attempts to match the parameters of the capability
	with compatible tags from the entry.
*/
func (capability Capability) CheckCompatability(entry cmdb.Entity) (bool, Capability) {
	var success bool

	// Check the capability preconditions are satisified, 1 from each group must be satisified
outer:
	for _, precon := range capability.Preconditions {
		preconSatisfied := false

		for _, sysTags := range entry.SysTags {
			// Label and DataType matches
			if precon.Label == sysTags.Label && precon.DataType == sysTags.DataType {
				for _, preconVal := range precon.Values {
					if utils.ArrayContains(preconVal, sysTags.Values) {
						preconSatisfied = true
						continue outer
					}
				}
			}
		}

		// Precondition not satisffied
		if !preconSatisfied {
			return false, capability
		}
	}

	// Check the command paramters
	for k, capParam := range capability.Command.Params {
		success, capability.Command.Params[k] = capParam.extractCompatibleParams(entry.SysTags)

		if !success {
			return false, capability
		}
	}

	return true, capability
}

/*
	extractCompatibleParams Determines if given a capability param {"Value": "","DataType": 1, "Default": ""}
	Is there any SysTags that can fulfil the Values
*/
func (capParam Param) extractCompatibleParams(entryTags []cmdb.EntityTag) (bool, Param) {
	// For each: {DataType.CMDB, DataType.IP}
	for _, pType := range capParam.DataType {
		// If the value is already set, or there's an available default, move on
		if capParam.Value != "" || capParam.Default != "" {

			if capParam.Value == "" {
				capParam.Value = capParam.Default
			}

			return true, capParam
		}

		// Skip empty tags that don't require input
		if pType == system.DataType_EMPTY {
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
