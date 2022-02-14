package execute

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/philcantcode/localmapper/capabilities/nmap"
	"github.com/philcantcode/localmapper/utils"
)

// Run commands from the console return: (output, error)
func Run(interpreter string, command string, args ...string) (interface{}, bool) {
	utils.Log("Attempting to run command: "+command, true)
	resultByte, err := exec.Command(command, args...).CombinedOutput()

	if err != nil {
		utils.ErrorHandle(fmt.Sprintf("Error returned running a command: %s", command), err, true)
		return nil, false
	}

	return interpret(string(resultByte), interpreter), true
}

func RunOnTop(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	utils.ErrorHandle("Couldn't run on top: "+command, err, false)
}

func interpret(result string, interpreter string) interface{} {
	switch interpreter {
	case "nmap:json":
		utils.Log("Interpreting nmap result", false)
		structuredResult := nmap.MakeStructured(result)
		utils.PrintLog(utils.PrettyPrintToStr(structuredResult))
		go nmap.InsertHosts(structuredResult)
		//go database.InsertMong(structuredResult)
		return structuredResult

	case "default":
		utils.Log("Interpreting default result", false)
		return result

	default:
		utils.Log("Interpreting defaulted result", false)
		return result
	}
}
