package capability

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/tools/nbtscan"
	"github.com/philcantcode/localmapper/tools/nmap"
	"github.com/philcantcode/localmapper/tools/searchsploit"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

var activeRoutines = 0
var maxRoutines = 8
var stopRoutines = false
var lifecycleCounter = 0

var queue = make(chan Lifecycle, maxRoutines)

type Lifecycle struct {
	ID           int
	capability   Capability
	resultBytes  []byte
	resultString string

	databaseInsertID string
	resultFilePath   string

	nmapRun    nmap.NmapRun
	nbtEntries []nbtscan.NBTScan
	exploitDBs []searchsploit.ExploitDB
}

func (lc *Lifecycle) SetCapability(capability Capability) {
	lc.capability = capability
	lc.ID = lifecycleCounter

	lifecycleCounter++

	system.Log(fmt.Sprintf("[Lifecycle Manager (%d)] Set Capability: %s", lc.ID, lc.capability.Label), true)
}

// 1 - Queue the excecution
func (lc *Lifecycle) Start() string {
	system.Log(fmt.Sprintf("[Queue]: %d/%d, adding: %s", len(queue), maxRoutines, lc.capability.Label), true)
	queue <- *lc

	return "TRACKING ID"
}

// 1 - Execution
func (lc *Lifecycle) execute() {
	system.Log(fmt.Sprintf("[Lifecycle Manager (%d)] Executing %s %v ...", lc.ID, lc.capability.Command.Program, lc.capability.ParamsToArray()), true)
	lc.resultBytes = execute(lc.capability.Command.Program, lc.capability.ParamsToArray())
}

// 2 - Interpret the results
func (lc *Lifecycle) interpret() {
	system.Log(fmt.Sprintf("[Lifecycle Manager (%d)] Interpreting ...", lc.ID), true)

	lc.resultString = string(lc.resultBytes)

	switch lc.capability.Interpreter {
	case system.Interpreter_NMAP:
		err := xml.Unmarshal(lc.resultBytes, &lc.nmapRun)
		system.Error("Couldn't unmarshal []bytes to NmapRun", err)
	case system.Interpreter_NBTSCAN:
		resultStrings := string(lc.resultBytes)
		resultArr := strings.Split(resultStrings, "\n")

		for _, line := range resultArr {
			lineArr := strings.Split(line, ",")

			if len(lineArr) >= 5 {
				nbtEntry := nbtscan.NBTScan{
					IP:          strings.TrimSpace(lineArr[0]),
					NetBIOSName: strings.TrimSpace(lineArr[1]),
					Server:      strings.TrimSpace(lineArr[2]),
					Username:    strings.TrimSpace(lineArr[3]),
					MAC:         strings.TrimSpace(lineArr[4]),
				}

				lc.nbtEntries = append(lc.nbtEntries, nbtEntry)
			}
		}
	case system.Interpreter_SEARCHSPLOIT:
		lc.exploitDBs = searchsploit.ExtractExploitDB(lc.resultBytes)
	default:
		system.Warning(fmt.Sprintf("No capability interpreter available for: %+v", lc.capability.Interpreter), true)
	}
}

// 3 - Records management & storage
func (lc *Lifecycle) recordsManagement() {
	system.Log(fmt.Sprintf("[Lifecycle Manager (%d)] Records Management ...", lc.ID), true)

	// NMAP
	if lc.capability.Interpreter == system.Interpreter_NMAP {
		lc.databaseInsertID = lc.nmapRun.Insert()

		// Write to nmap output directory
		lc.resultFilePath = fmt.Sprintf("%s/%s.txt", system.GetConfig("nmap-results-dir"), lc.databaseInsertID)
		utils.CreateAndWriteFile(lc.resultFilePath, string(lc.resultBytes))

		// Extract entities
		for _, entity := range lc.nmapRun.ExtractEntities() {
			entity.UpdateOrInsert()
		}
	}

	if lc.capability.Interpreter == system.Interpreter_NBTSCAN {
		for _, nbtEntry := range lc.nbtEntries {
			nbtEntry.Insert()

			// Extract entities
			for _, entity := range nbtEntry.ExtractEntities() {
				entity.UpdateOrInsert()
			}
		}
	}

	if lc.capability.Interpreter == system.Interpreter_SEARCHSPLOIT {
		for _, exploit := range lc.exploitDBs {

			// Don't input where the exact search has been done before
			if len(searchsploit.Select(bson.M{"search": exploit.Search}, bson.M{})) != 0 {
				searchsploit.Delete(bson.M{"search": exploit.Search})
			}

			exploit.Insert()
		}
	}
}

// 4 - Other add-ons
func (lc *Lifecycle) capabilityHooks() {
	system.Log(fmt.Sprintf("[Lifecycle Manager (%d)] CapabilityHooks ...", lc.ID), true)

	if lc.capability.Interpreter == system.Interpreter_NMAP {
		// Run results through searchsploit
		searchsploit := SELECT_Capability(bson.M{"cci": "cci:searchsploit:nmap:json"}, bson.M{})

		if len(searchsploit) != 1 {
			system.Warning("Couldn't find cci:searchsploit:nmap:json", false)
			return
		}

		// Search for vulnerbilities
		searchsploit[0].Command.Params[0].Value = lc.resultFilePath
		searchSploitLC := Lifecycle{capability: searchsploit[0]}
		searchSploitLC.Start()
	}
}

/*
	ProcessCapabilityQueue should be called once at the start
	of the application to start the capability processing.
*/
func InitCapabilityLifecycleManager() {
	for {
		if stopRoutines {
			break
		}

		if activeRoutines < maxRoutines {
			go func() {
				activeRoutines++

				lc := <-queue

				lc.execute()
				lc.interpret()
				lc.recordsManagement()
				lc.capabilityHooks()

				system.Log(fmt.Sprintf("[Capability Queue]: %d/%d", len(queue), maxRoutines), true)

				activeRoutines--
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
	stopRoutines = true
	system.Log("Stopping queue, currently processing jobs will continue", false)
}
