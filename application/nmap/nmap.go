package nmap

import (
	"encoding/xml"
	"fmt"
	"os/exec"

	"github.com/philcantcode/localmapper/adapters/definitions"
	"github.com/philcantcode/localmapper/utils"
)

func RunNmapCommand(capability definitions.Capability) definitions.NmapRun {

	prog := capability.Command.Program
	params := definitions.ParamsToArray(capability.Command.Params)

	utils.Log(fmt.Sprintf("Attempting to run Nmap Command: %s > %v", prog, params), true)
	resultByte, err := exec.Command(prog, params...).CombinedOutput()
	utils.ErrorFatal(fmt.Sprintf("Error returned running a command: %s > %v", prog, params), err)

	return interpret(resultByte)
}

func interpret(result []byte) definitions.NmapRun {
	utils.Log("Converting from []byte to NmapRun struct", false)

	var nmapRun definitions.NmapRun
	err := xml.Unmarshal(result, &nmapRun)

	utils.ErrorFatal("Couldn't unmarshal result from Nmap console", err)

	return nmapRun
}
