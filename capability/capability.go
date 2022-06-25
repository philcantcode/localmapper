package capability

import (
	"fmt"
	"time"

	"github.com/philcantcode/localmapper/capability/acccheck"
	"github.com/philcantcode/localmapper/capability/local"
	"github.com/philcantcode/localmapper/capability/nbtscan"
	"github.com/philcantcode/localmapper/capability/nmap"
	"github.com/philcantcode/localmapper/capability/searchsploit"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

var currentRoutines = 0
var maxRoutines = 8
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
				system.Log(fmt.Sprintf("[Capability Queue]: %d/%d", len(queue), maxRoutines), true)

				currentRoutines--
			}()
		}

		time.Sleep(3000)
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
	case system.Interpreter_NMAP:
		resultBytes := local.Execute(capability.Command.Program, ParamsToArray(capability.Command.Params))
		nmapRun := nmap.ProcessResults(resultBytes)
		nmapRun.ConvertToEntry()
		resID := nmapRun.StoreResults()
		path := nmap.WriteResultToDisk(resultBytes, resID)

		// Run results through searchsploit
		searchsploit := SELECT_Capability(bson.M{"cci": "cci:searchsploit:nmap:json"}, bson.M{})

		if len(searchsploit) != 1 {
			system.Force("Couldn't find cci:searchsploit:nmap:json", false)
			return
		}

		searchsploit[0].Command.Params[0].Value = path
		go searchsploit[0].QueueCapability()
	case system.Interpreter_UNIVERSAL:
		res := local.Execute(capability.Command.Program, ParamsToArray(capability.Command.Params))
		system.Log("UNIVERSAL RESULT : "+string(res), true)
	case system.Interpreter_ACCCHECK:
		resultBytes := local.Execute(capability.Command.Program, ParamsToArray(capability.Command.Params))
		result := acccheck.ProcessResults(resultBytes)
		acccheck.StoreResults(result)
	case system.Interpreter_NBTSCAN:
		resultBytes := local.Execute(capability.Command.Program, ParamsToArray(capability.Command.Params))
		nbtScan := nbtscan.ProcessResults(resultBytes)
		nbtscan.ConvertToEntry(nbtScan)
		nbtscan.StoreResults(nbtScan)
	case system.Interpreter_SEARCHSPLOIT:
		resultBytes := local.Execute(capability.Command.Program, ParamsToArray(capability.Command.Params))
		exploitDB := searchsploit.ProcessResults(resultBytes)

		for _, exp := range exploitDB {
			exp.StoreResults()
		}
	default:
		system.Force(fmt.Sprintf("No capability interpreter available for: %+v", capability), true)
	}
}

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
