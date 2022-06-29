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

var queue = make(chan *Lifecycle, maxRoutines)
var managementStore = []*Lifecycle{}

type Status string

const (
	Status_Created           Status = "CREATED"
	Status_Waiting           Status = "WAITING"
	Status_Executing         Status = "EXECUTING"
	Status_Interpreting      Status = "INTERPRETING"
	Status_RecordsManagement Status = "RECORDS_MANAGEMENT"
	Status_CapabilityHooks   Status = "CAPABILITY_HOOKS"
	Status_Done              Status = "DONE"
)

/*
	Lifecycle keeps track of the capability execution
	and interpretation through its lifecycle
*/
type Lifecycle struct {
	Tracking     Tracking
	Capability   Capability
	ResultBytes  []byte
	ResultString string

	// Recrods management
	DatabaseInsertID string
	ResultFilePath   string

	// Structures that may be used
	nmapRun    nmap.NmapRun
	nbtEntries []nbtscan.NBTScan
	exploitDBs []searchsploit.ExploitDB
}

type Tracking struct {
	ID                   int
	Command              string
	RuntimeStart         time.Time
	RuntimeEnd           time.Time
	RuntimeDuration      time.Duration
	RuntimeDurationPrint string
	ExecStart            time.Time
	ExecEnd              time.Time
	ExecDuration         time.Duration
	ExecDurationPrint    string
	Status               Status
}

func (lc *Lifecycle) SetCapability(capability Capability) {
	lc.Tracking.Status = Status_Created
	lc.Capability = capability
	lc.Tracking.ID = lifecycleCounter
	managementStore = append(managementStore, lc)

	lifecycleCounter++

	system.Log(fmt.Sprintf("[Lifecycle Manager (%d)] Set Capability: %s", lc.Tracking.ID, lc.Capability.Label), false)
}

// 1 - Queue the excecution
func (lc *Lifecycle) Start() string {
	lc.Tracking.Status = Status_Waiting
	lc.Tracking.RuntimeStart = time.Now()
	lc.Tracking.Command = fmt.Sprintf("%s %v", lc.Capability.Command.Program, lc.Capability.ParamsToArray())

	system.Log(fmt.Sprintf("[Queue]: %d/%d, adding: %s", len(queue), maxRoutines, lc.Capability.Label), true)
	queue <- lc

	return "TRACKING ID"
}

// 1 - Execution
func (lc *Lifecycle) execute() {
	lc.Tracking.Status = Status_Executing
	system.Log(fmt.Sprintf("[Lifecycle Manager (%d)] Executing %s %v ...", lc.Tracking.ID, lc.Capability.Command.Program, lc.Capability.ParamsToArray()), true)
	lc.ResultBytes = execute(lc.Capability.Command.Program, lc.Capability.ParamsToArray())
}

// 2 - Interpret the results
func (lc *Lifecycle) interpret() {
	lc.Tracking.Status = Status_Interpreting
	system.Log(fmt.Sprintf("[Lifecycle Manager (%d)] Interpreting ...", lc.Tracking.ID), false)

	lc.ResultString = string(lc.ResultBytes)

	switch lc.Capability.Interpreter {
	case system.Interpreter_NMAP:
		err := xml.Unmarshal(lc.ResultBytes, &lc.nmapRun)
		system.Error("Couldn't unmarshal []bytes to NmapRun", err)
	case system.Interpreter_NBTSCAN:
		resultStrings := string(lc.ResultBytes)
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
		lc.exploitDBs = searchsploit.ExtractExploitDB(lc.ResultBytes)
	default:
		system.Warning(fmt.Sprintf("No capability interpreter available for: %+v", lc.Capability.Interpreter), true)
	}
}

// 3 - Records management & storage
func (lc *Lifecycle) recordsManagement() {
	lc.Tracking.Status = Status_RecordsManagement
	system.Log(fmt.Sprintf("[Lifecycle Manager (%d)] Records Management ...", lc.Tracking.ID), false)

	// NMAP
	if lc.Capability.Interpreter == system.Interpreter_NMAP {
		lc.DatabaseInsertID = lc.nmapRun.Insert()

		// Write to nmap output directory
		lc.ResultFilePath = fmt.Sprintf("%s/%s.txt", system.GetConfig("nmap-results-dir"), lc.DatabaseInsertID)
		utils.CreateAndWriteFile(lc.ResultFilePath, string(lc.ResultBytes))

		// Extract entities
		for _, entity := range lc.nmapRun.ExtractEntities() {
			entity.UpdateOrInsert()
		}
	}

	if lc.Capability.Interpreter == system.Interpreter_NBTSCAN {
		for _, nbtEntry := range lc.nbtEntries {
			nbtEntry.Insert()

			// Extract entities
			for _, entity := range nbtEntry.ExtractEntities() {
				entity.UpdateOrInsert()
			}
		}
	}

	if lc.Capability.Interpreter == system.Interpreter_SEARCHSPLOIT {
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
	lc.Tracking.Status = Status_CapabilityHooks
	system.Log(fmt.Sprintf("[Lifecycle Manager (%d)] CapabilityHooks ...", lc.Tracking.ID), false)

	if lc.Capability.Interpreter == system.Interpreter_NMAP {
		// Run results through searchsploit
		searchsploit := SELECT_Capability(bson.M{"cci": "cci:searchsploit:nmap:json"}, bson.M{})

		if len(searchsploit) != 1 {
			system.Warning("Couldn't find cci:searchsploit:nmap:json", false)
			return
		}

		// Search for vulnerbilities
		searchsploit[0].Command.Params[0].Value = lc.ResultFilePath
		searchSploitLC := Lifecycle{Capability: searchsploit[0]}
		go searchSploitLC.Start()
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
				lc.Tracking.ExecStart = time.Now()

				lc.execute()
				lc.interpret()
				lc.recordsManagement()
				lc.capabilityHooks()

				lc.Tracking.RuntimeEnd = time.Now()
				lc.Tracking.ExecEnd = time.Now()

				lc.Tracking.RuntimeDuration = lc.Tracking.RuntimeEnd.Sub(lc.Tracking.RuntimeStart)
				lc.Tracking.ExecDuration = lc.Tracking.ExecEnd.Sub(lc.Tracking.ExecStart)

				lc.Tracking.RuntimeDurationPrint = utils.FormatDuration(lc.Tracking.RuntimeDuration)
				lc.Tracking.ExecDurationPrint = utils.FormatDuration(lc.Tracking.ExecDuration)
				lc.Tracking.Status = Status_Done
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
