package nmap

import (
	"encoding/xml"
	"fmt"
	"os/exec"

	"github.com/philcantcode/localmapper/adapters/blueprint"
	"github.com/philcantcode/localmapper/utils"
)

func RunNmapCommand(capability blueprint.Capability) blueprint.NmapRun {
	utils.Log(fmt.Sprintf("Attempting to run Nmap Command: %v", capability), true)
	resultByte, err := exec.Command(capability.Params[0], capability.Params[1:]...).CombinedOutput()

	utils.ErrorFatal(fmt.Sprintf("Error returned running a command: %v", capability), err)

	return interpret(resultByte)
}

func interpret(result []byte) blueprint.NmapRun {
	utils.Log("Converting from []byte to NmapRun struct", true)

	var nmapRun blueprint.NmapRun
	err := xml.Unmarshal(result, &nmapRun)

	utils.ErrorFatal("Couldn't unmarshal result from Nmap console", err)

	return nmapRun
}
