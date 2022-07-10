package cookbook

import (
	"fmt"
	"time"

	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/utils"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type scheduleTracking struct {
	CCBI  string
	Index int // Schedule index
	Cron  *cron.Cron
}

var tracking []Tracking

type Schedule struct {
	Label         string
	Desc          string
	Delay         time.Duration   // Time delay between runs
	RescanDelay   time.Duration   // Time between re-scanning entries once they've already been scanned
	TargetGroups  []cmdb.CMDBType // ALL_PENDING, ALL_CMDB etc
	TargetDevices []string        // Individual device IDs
	ExclusionList []Exclusion     // Devices to not be scanned
	Tracking      Tracking        // Metadata to for scheduling
}

/*
	Data used to keep track of the schedule processing
*/
type Tracking struct {
	ID           int
	IsRunning    bool                 // Keeps track of whether the cookbook is running currently
	EntryHistory map[string]time.Time // Keeps an [ID:Time] map of the last time each Entry ran
}

/*
	Definition of an IP range to be excluded
*/
type Exclusion struct {
	Value    string
	DataType system.DataType
}

/*
	Assumes none currently running & tracking array empty
*/
func InitialiseAllSchedules() {
	books := SELECT_Cookbook(bson.M{}, bson.M{})
	scheduleIdx := 0

	// Initialise the tracking IDs
	for bidx, book := range books {
		for sidx, _ := range book.Schedule {
			trackDetails := Tracking{ID: scheduleIdx, IsRunning: false}
			tracking = append(tracking, trackDetails)
			books[bidx].Schedule[sidx].Tracking = trackDetails
			scheduleIdx++
		}
	}

	// Loop over and initialise jobs if not running
	for _, book := range books {
		for sidx, schedule := range book.Schedule {
			system.Log("Scheduling new cron job: "+schedule.Label, true)
			cronny := cron.New()

			/*
				Make sure you create local instances of the variables
				to prevent race conditions, otherwise every cron job
				will process the last variable.

				https://stackoverflow.com/questions/28954869/cannot-assign-variable-to-anonymous-func-in-for-loop
			*/
			jBook := book
			jSchedule := schedule
			cronny.AddJob("@every "+schedule.Delay.String(), cron.FuncJob(func() {

				if tracking[sidx].IsRunning {
					system.Log(fmt.Sprintf("Schedule [%s] already running, skipping.", jSchedule.Label), true)
					return
				}

				tracking[sidx].IsRunning = true

				for _, targID := range jSchedule.TargetDevices {
					targ := cmdb.SELECT_Entities_Joined(bson.M{"_id": system.EncodeID(targID)}, bson.M{})

					if len(targ) != 1 {
						system.Warning(fmt.Sprintf("Incorrect number of targets returned (%d) for %s", len(targ), targID), true)
					}

					// Check whether suitable time has past between last entry scan
					var hasTimePassed bool
					jSchedule, hasTimePassed = hasSuitableTimePassed(jSchedule, targ[0].ID.Hex())
					isExcluded := isExcluded(targ[0], jSchedule.ExclusionList)

					if !isExcluded && hasTimePassed {
						system.Log(fmt.Sprintf("Starting Schedule: %s (ccbi: %s) > %s", jSchedule.Label, jBook.CCBI, targ[0].Label), true)
						jBook.RunBookOnEntity(targ[0].ID)
					} else {
						system.Log(fmt.Sprintf("Skipping %s because [ON EXCLUSION LIST: %t] [SUITABLE TIME PASSED: %t]\n", targ[0].Label, isExcluded, hasTimePassed), false)
					}
				}

				for _, targGroup := range jSchedule.TargetGroups {
					for _, entry := range cmdb.SELECT_Entities_Joined(bson.M{"cmdbtype": targGroup}, bson.M{}) {

						// Check whether suitable time has past between last entry scan
						var hasTimePassed bool
						jSchedule, hasTimePassed = hasSuitableTimePassed(jSchedule, entry.ID.Hex())
						isExcluded := isExcluded(entry, jSchedule.ExclusionList)

						if !isExcluded && hasTimePassed {
							system.Log(fmt.Sprintf("Starting Schedule: %s (ccbi: %s) > %s", jSchedule.Label, jBook.CCBI, entry.Label), true)
							jBook.RunBookOnEntity(entry.ID)
						} else {
							system.Log(fmt.Sprintf("Skipping %s because [ON EXCLUSION LIST: %t] [SUITABLE TIME PASSED: %t]", entry.Label, isExcluded, hasTimePassed), false)
						}
					}
				}

				tracking[sidx].IsRunning = false
			}))

			cronny.Start()
		}
	}
}

/*
	isExcluded searches the SysTags for any exclusion values
*/
func isExcluded(entry cmdb.Entity, exclusions []Exclusion) bool {
	for _, tag := range entry.SysTags {
		for _, excl := range exclusions {
			// If the DataTypes match and exclue value in tags array
			if tag.DataType == excl.DataType && utils.ArrayContains(excl.Value, tag.Values) {
				return true
			}
		}
	}

	return false
}

/*
	hasSuitableTimePassed returns a bool whether suitable amount
	of time has passed between this and the last scan
*/
func hasSuitableTimePassed(s Schedule, eID string) (Schedule, bool) {

	// initialise map
	if len(s.Tracking.EntryHistory) == 0 {
		s.Tracking.EntryHistory = make(map[string]time.Time)
	}

	if lastScan, ok := s.Tracking.EntryHistory[eID]; ok {

		if time.Now().After(lastScan.Add(s.RescanDelay)) {
			s.Tracking.EntryHistory[eID] = time.Now()
			return s, true
		}

		return s, false
	}

	s.Tracking.EntryHistory[eID] = time.Now()
	// fmt.Printf("%+v\n", s.Tracking.EntryHistory)
	return s, true
}
