package capability

import (
	"time"

	"github.com/philcantcode/localmapper/utils"
)

type TimeGraph struct {
	Keys   []string
	Values []int
}

type JobTypeGraph struct {
	Value int
	Name  string
}

func CalcJobsTimeGraph() TimeGraph {
	graph := TimeGraph{Keys: []string{}, Values: []int{}}

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
