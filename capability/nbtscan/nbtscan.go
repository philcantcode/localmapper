package nbtscan

import (
	"strings"
)

type NBTScan struct {
	IP          string
	NetBIOSName string
	Server      string
	Username    string
	MAC         string
}

func Interpret(resultBytes []byte) []NBTScan {

	results := string(resultBytes)
	resultArr := strings.Split(results, "\n")
	nbtResults := []NBTScan{}

	for _, line := range resultArr {
		lineArr := strings.Split(line, ",")

		if len(lineArr) >= 5 {
			res := NBTScan{
				IP:          strings.TrimSpace(lineArr[0]),
				NetBIOSName: strings.TrimSpace(lineArr[1]),
				Server:      strings.TrimSpace(lineArr[2]),
				Username:    strings.TrimSpace(lineArr[3]),
				MAC:         strings.TrimSpace(lineArr[4]),
			}

			nbtResults = append(nbtResults, res)
		}
	}

	return nbtResults
}
