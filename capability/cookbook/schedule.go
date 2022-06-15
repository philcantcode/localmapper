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
	Duration      time.Duration   // Time between runs
	TargetGroups  []cmdb.CMDBType // ALL_PENDING, ALL_CMDB etc
	TargetDevices []string        // Individual device IDs
	ExclusionList []Exclusion     // Devices to not be scanned
	Tracking      Tracking        // Metadata to for scheduling
}

/*
	Data used to keep track of the schedule processing
*/
type Tracking struct {
	ID        int
	IsRunning bool
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
			system.Log("Scheduling new task: "+schedule.Label, true)
			cronny := cron.New()

			cronny.AddJob("@every "+schedule.Duration.String(), cron.FuncJob(func() {
				if tracking[sidx].IsRunning {
					system.Log(fmt.Sprintf("Schedule [%s] already running, skipping.", schedule.Label), true)
					return
				}

				tracking[sidx].IsRunning = true

				system.Log(fmt.Sprintf("Starting scheduled task %s", schedule.Label), true)

				for _, targID := range schedule.TargetDevices {
					targ := cmdb.SELECT_ENTRY_Joined(bson.M{"_id": system.EncodeID(targID)}, bson.M{})

					if len(targ) != 1 {
						system.Force(fmt.Sprintf("Incorrect number of targets returned (%d) for %s", len(targ), targID), true)
					}

					if !isExcluded(targ[0], schedule.ExclusionList) {
						ExecuteCookbook(book, targ[0].ID)
					}
				}

				for _, targGroup := range schedule.TargetGroups {
					for _, e := range cmdb.SELECT_ENTRY_Joined(bson.M{"cmdbtype": targGroup}, bson.M{}) {
						if !isExcluded(e, schedule.ExclusionList) {
							ExecuteCookbook(book, e.ID)
						}
					}
				}

				system.Log(fmt.Sprintf("Finishing scheduled task %s", schedule.Label), true)
				tracking[sidx].IsRunning = false
				time.Sleep(3000)
			}))

			cronny.Start()
		}
	}
}

/*
	isExcluded searches the SysTags for any exclusion values
*/
func isExcluded(entry cmdb.Entry, exclusions []Exclusion) bool {
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
