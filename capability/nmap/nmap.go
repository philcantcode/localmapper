package nmap

import (
	"encoding/xml"
	"fmt"
	"os/exec"

	"github.com/philcantcode/localmapper/utils"
)

// Execute takes a list of parameters to execute against NMAP
func Execute(params []string) NmapRun {
	utils.Log(fmt.Sprintf("Attempting to run Nmap Command: %s > %v", "nmap", params), true)
	resultByte, err := exec.Command("nmap", params...).CombinedOutput()
	utils.ErrorFatal(fmt.Sprintf("Error returned running a command: %s > %v", "nmap", params), err)

	return interpret(resultByte)
}

func interpret(result []byte) NmapRun {
	utils.Log("Converting from []byte to NmapRun struct", false)

	var nmapRun NmapRun
	err := xml.Unmarshal(result, &nmapRun)

	utils.ErrorFatal("Couldn't unmarshal result from Nmap console", err)

	return nmapRun
}
