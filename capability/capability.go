package capability

import (
	"fmt"
	"time"

	"github.com/philcantcode/localmapper/capability/acccheck"
	"github.com/philcantcode/localmapper/capability/local"
	"github.com/philcantcode/localmapper/capability/nbtscan"
	"github.com/philcantcode/localmapper/capability/nmap"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
)

var currentRoutines = 0
var maxRoutines = 50
var queue = make(chan Capability, maxRoutines)
var stopCapability = false

func (capability Capability) QueueCapability() {
	system.Log(fmt.Sprintf("[Capability Queue]: %d/%d, adding: %s", len(queue), maxRoutines, capability.Label), true)
	queue <- capability
}

/*
	ProcessCapabilityQueue should be called once at the start
	of the application to start the capability processing.
*/
func ProcessCapabilityQueue() {
	for {
		if stopCapability {
			break
		}

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

/*
	StopCapabilityQueue breaks out of the processing loop
	so you must call ProcessCapabilityQueue again to resume.

	It allows currently running capabilities to finish.
*/
func StopCapabilityQueue() {
	stopCapability = true
	system.Log("Stopping queue, currently processing jobs will continue", false)
}

/*
	executeCapability executes a capability then attempts to process the result.
	Each result type must implement a "ProcessResults()" function to pack the
	byte[] into a struct and then a "StoreResults()" function to store it in the
	relevant databases.
*/
func executeCapability(capability Capability) {
	system.Log(fmt.Sprintf("Executing Capability: %s", capability.Label), true)

	switch capability.Interpreter {
	case system.NMAP:
		resultBytes := local.Execute(capability.Command.Program, ParamsToArray(capability.Command.Params))
		nmapRun := nmap.ProcessResults(resultBytes)
		nmapRun.ConvertToEntry()
		nmapRun.StoreResults()
	case system.UNIVERSAL:
		local.Execute(capability.Command.Program, ParamsToArray(capability.Command.Params))
	case system.ACCCHECK:
		resultBytes := local.Execute(capability.Command.Program, ParamsToArray(capability.Command.Params))
		result := acccheck.ProcessResults(resultBytes)
		acccheck.StoreResults(result)
	case system.NBTSCAN:
		resultBytes := local.Execute(capability.Command.Program, ParamsToArray(capability.Command.Params))
		nbtScan := nbtscan.ProcessResults(resultBytes)
		nbtscan.ConvertToEntry(nbtScan)
		nbtscan.StoreResults(nbtScan)
	default:
		system.Force(fmt.Sprintf("No capability interpreter available for: %+v", capability), true)
	}
}

/*
	ExtractCompabileTags takes a given entry and attempts to match the parameters of the capability
	with compatible tags from the entry.
*/
func (capability Capability) ExtractCompabileTags(entry cmdb.Entry) (bool, Capability) {
	var success bool

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
func (capParam Param) extractCompatibleParams(entryTags []cmdb.EntryTag) (bool, Param) {
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
