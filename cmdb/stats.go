package cmdb

import (
	"strconv"
	"time"

	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/utils"
)

type IdentityConfidence struct {
	IP       int
	MAC      int
	HostName int
	Vendor   int
	OS       int
	DateSeen int

	Total   int
	Average int
}

type TimeGraph struct {
	Keys   []string
	Values []int
}

type CMDBEntityStats struct {
	Confidence IdentityConfidence
	TimeGraph  TimeGraph
}

func (entry *Entity) GetStats() CMDBEntityStats {
	stats := CMDBEntityStats{
		Confidence: entry.calcIdentityConfidenceScore(),
		TimeGraph:  entry.calcTimeGraph(),
	}

	return stats
}

/*
	CalcIdentityConfidenceScore returns an structure of confidences for various
	values known about the device.
*/
func (entry *Entity) calcIdentityConfidenceScore() IdentityConfidence {
	result := IdentityConfidence{IP: 0, MAC: 0, HostName: 0, Vendor: 0, OS: 0, DateSeen: 0, Average: 0}

	// IP4 weighted 80%, divided by the number of past changes
	ipTag, found, _ := entry.FindSysTag("IP")

	if found {
		result.IP += (80 / len(ipTag.Values))
	}

	// IP6 weigthed 20%, divided by the number of past changes
	ip6Tag, found, _ := entry.FindSysTag("IP6")

	if found {
		result.IP += (20 / len(ip6Tag.Values))
	}

	// MAC weighted 80%, divided by the number of past changes
	macTag, found, _ := entry.FindSysTag("MAC")

	if found {
		result.MAC += (80 / len(macTag.Values))
	}

	// MAC6 weigthed 20%, divided by the number of past changes
	mac6Tag, found, _ := entry.FindSysTag("MAC6")

	if found {
		result.MAC += (20 / len(mac6Tag.Values))
	}

	// Host Name weigthed 100%, divided by the number of past changes
	hostNameTag, found, _ := entry.FindSysTag("HostName")

	if found {
		result.HostName = (100 / len(hostNameTag.Values))
	}

	_, hasMACVendor, _ := entry.FindSysTag("MACVendor")
	_, hasOSVendor, _ := entry.FindSysTag("OSVendor")
	vendorACC, hasVendorACC, _ := entry.FindSysTag("OSAccuracy")

	if hasVendorACC && hasOSVendor {
		vAccInt, err := strconv.Atoi(vendorACC.Values[0])
		system.Fatal("Couldn't convert CMDBType to int", err)

		result.Vendor += (vAccInt / 2)
	}

	if hasMACVendor {
		result.Vendor += 50
	}

	_, found, _ = entry.FindSysTag("OS")

	if found {
		result.OS = 100
	}

	// Assign values based on the number of dates seen
	size := len(entry.DateSeen)
	if size == 0 {
		result.DateSeen = 0
	} else if size > 0 && size <= 5 {
		result.DateSeen = 33
	} else if size > 5 && size <= 10 {
		result.DateSeen = 50
	} else if size > 10 && size <= 50 {
		result.DateSeen = 75
	} else if size > 50 {
		result.DateSeen = 100
	}

	// Calcualte an average
	result.Total = result.IP + result.MAC + result.OS + result.Vendor + result.DateSeen + result.HostName
	result.Average = result.Total / 6

	return result
}

func (entry *Entity) calcTimeGraph() TimeGraph {
	graph := TimeGraph{Keys: []string{}, Values: []int{}}

	windowMinute, err := strconv.Atoi(system.GetConfig("date-seen-graph-mins-val"))
	system.Fatal("Couldn't convert date-seen-graph-mins-val settings to int", err)

	window := time.Duration(int64(time.Minute) * int64(windowMinute))
	nowDT, _ := time.Parse(utils.DTF_DateTime, entry.DateSeen[0])
	endDT, _ := time.Parse(utils.DTF_DateTime, entry.DateSeen[len(entry.DateSeen)-1])
	processedBlocks := 1
	entriesInBlock := 0
	timeBlockCount := 0

	// Push the first value of DateTime to get things started
	graph.Keys = append(graph.Keys, nowDT.Format(utils.DTF_DateTime))

	for nowDT.Before(endDT) {
		nowDT = nowDT.Add(window) // Add 1 duration step to it

		for _, dt := range entry.DateSeen[processedBlocks:] {
			subDT, _ := time.Parse(utils.DTF_DateTime, dt)

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
