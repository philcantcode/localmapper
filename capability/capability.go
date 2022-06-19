package capability

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/philcantcode/localmapper/capability/acccheck"
	"github.com/philcantcode/localmapper/capability/local"
	"github.com/philcantcode/localmapper/capability/nbtscan"
	"github.com/philcantcode/localmapper/capability/nmap"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/interpreter"
	"github.com/philcantcode/localmapper/system"
)

var currentRoutines = 0
var maxRoutines = 30
var queue = make(chan Capability, maxRoutines)

func QueueCapability(capability Capability) {
	system.Log(fmt.Sprintf("[Current Queue Size]: %d, adding: %s", len(queue), capability.Label), true)
	queue <- capability
}

func ProcessCapabilityQueue() {
	for {
		if currentRoutines < maxRoutines {
			go func() {
				currentRoutines++
				cap := <-queue
				executeCapability(cap)
				currentRoutines--
			}()
		}

		time.Sleep(1000)
	}
}

func executeCapability(capability Capability) []byte {
	system.Log(fmt.Sprintf("Executing Capability: %s", capability.Label), true)

	switch capability.Interpreter {
	case system.NMAP:
		resultBytes := local.Execute(capability.Command.Program, ParamsToArray(capability.Command.Params))

		nmapRun := nmap.Interpret(resultBytes)
		nmap.INSERT_Nmap(nmapRun)

		result, err := json.Marshal(nmapRun)
		system.Error("Couldn't marshal nmaprun", err)

		return result
	case system.UNIVERSAL:
		resultBytes := local.Execute(capability.Command.Program, ParamsToArray(capability.Command.Params))

		return interpreter.UniversalExec(resultBytes)
	case system.ACCCHECK:
		resultBytes := local.Execute(capability.Command.Program, ParamsToArray(capability.Command.Params))
		result := acccheck.Interpret(resultBytes)

		return result
	case system.NBTSCAN:
		resultBytes := local.Execute(capability.Command.Program, ParamsToArray(capability.Command.Params))
		results := nbtscan.Interpret(resultBytes)

		for _, res := range results {
			nbtscan.INSERT_NbtScan(res)
		}

		return nil
	default:
		system.Force(fmt.Sprintf("No capability interpreter available for: %+v", capability), true)
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
