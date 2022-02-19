package execute

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/philcantcode/localmapper/apps/nmap"
	"github.com/philcantcode/localmapper/utils"
)

// Run commands from the console return
func Run(capability Capability) interface{} {
	utils.Log(fmt.Sprintf("Attempting to run command: %v", capability), true)
	resultByte, err := exec.Command(capability.Params[0], capability.Params[1:]...).CombinedOutput()

	utils.ErrorFatal(fmt.Sprintf("Error returned running a command: %v", capability), err)

	return interpret(string(resultByte), capability.Interpreter)
}

// RunToStdOut prints the results to STDOUT
func RunToStdOut(capability Capability) {
	console := exec.Command(capability.Params[0], capability.Params[1:]...)
	console.Stdout = os.Stdout
	console.Stderr = os.Stderr

	err := console.Run()
	utils.ErrorFatal(fmt.Sprintf("Couldn't RunOnTop: %v", capability), err)
}

func interpret(result string, interpreter string) interface{} {
	switch interpreter {
	case "nmap:json":
		utils.Log("Interpreting nmap result", false)
		structuredResult := nmap.MakeStructured(result)
		utils.PrintLog(utils.PrettyPrintToStr(structuredResult))
		go nmap.SqliteInsertHosts(structuredResult)
		go nmap.MongoInsert(structuredResult)
		return structuredResult

	case "default":
		utils.Log("Interpreting default result", false)
		return result

	default:
		utils.Log("Interpreting defaulted result", false)
		return result
	}
}
