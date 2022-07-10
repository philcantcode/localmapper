package capability

import (
	"time"

	"github.com/philcantcode/localmapper/utils"
)

type JobStats struct {
	Running       int
	Waiting       int
	TimeGraph     JobTimeGraph
	TypeGraph     []JobTypeGraph
	LifecycleJobs []*Lifecycle
}

type JobTimeGraph struct {
	Keys   []string
	Values []int
}

type JobTypeGraph struct {
	Value int
	Name  string
}

/*
	Retruns stats about running jobs, including:

	1. Timings of running jobs
	2. List of all job categories (past and current)
	3. All lifecycle management jobs (e.g., with status)
	4. Running jobs
	5. Waiting jobs
*/
func GetJobStats() JobStats {

	res := JobStats{
		TimeGraph:     generateJobTimingGraph(),
		TypeGraph:     generateJobTypeGraph(),
		LifecycleJobs: getAllLifecycleTasks(),
	}

	for _, job := range managementStore {
		if job.Tracking.Status != Status_Done && job.Tracking.Status != Status_Waiting {
			res.Running++
		}

		if job.Tracking.Status != Status_Done && job.Tracking.Status == Status_Waiting {
			res.Waiting++
		}
	}

	return res
}

func generateJobTypeGraph() []JobTypeGraph {
	jobTypeGraph := []JobTypeGraph{}

	// Get a running value for duration
	for _, v := range managementStore {
		name := v.Capability.Label
		found := false

		for idx, k := range jobTypeGraph {
			if k.Name == name {
				jobTypeGraph[idx].Value++
				found = true
				break
			}
		}

		if !found {
			jobTypeGraph = append(jobTypeGraph, JobTypeGraph{Name: name, Value: 1})
		}
	}

	return jobTypeGraph
}

func generateJobTimingGraph() JobTimeGraph {
	graph := JobTimeGraph{Keys: []string{}, Values: []int{}}

	if len(managementStore) == 0 {
		return graph
	}

	windowMinute := 10

	window := time.Duration(int64(time.Second) * int64(windowMinute))
	nowDT := managementStore[0].Tracking.RuntimeStart
	endDT := managementStore[len(managementStore)-1].Tracking.RuntimeStart

	processedBlocks := 1
	entriesInBlock := 0
	timeBlockCount := 0

	// Push the first value of DateTime to get things started
	graph.Keys = append(graph.Keys, nowDT.Format(utils.DTF_DateTime))

	for nowDT.Before(endDT) {
		nowDT = nowDT.Add(window) // Add 1 duration step to it

		for _, dt := range managementStore[processedBlocks:] {
			subDT := dt.Tracking.RuntimeStart

			if subDT.Before(nowDT) {
				entriesInBlock++
			} else {
				break
			}

			processedBlocks++
		}

		graph.Keys = append(graph.Keys, nowDT.Format(utils.DTF_DateTime))
		graph.Values = append(graph.Values, entriesInBlock)

		entriesInBlock = 0
		timeBlockCount++
	}

	return graph
}

/*
	Calculates the running time
*/
func getAllLifecycleTasks() []*Lifecycle {
	// Get a running value for duration
	for i := range managementStore {
		if managementStore[i].Tracking.Status != Status_Done {
			managementStore[i].Tracking.RuntimeDuration = time.Now().Sub(managementStore[i].Tracking.RuntimeStart)
			managementStore[i].Tracking.ExecDuration = time.Now().Sub(managementStore[i].Tracking.ExecStart)

			managementStore[i].Tracking.RuntimeDurationPrint = utils.FormatDuration(managementStore[i].Tracking.RuntimeDuration)
			managementStore[i].Tracking.ExecDurationPrint = utils.FormatDuration(managementStore[i].Tracking.ExecDuration)
		}
	}

	return managementStore
}
