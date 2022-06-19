package system

import (
	"fmt"
	"os/exec"
)

func SetupTools() {

	// Acc check, set permissions then execute
	execute("chmod +x ./res/apps/acccheck/acccheck.pl")
	execute("./res/apps/acccheck/acccheck.pl")
}

func UnpackWordLists() {

}

func execute(command string) []byte {
	Log(fmt.Sprintf("Attempting system.Execute: %s", command), true)

	resultByte, err := exec.Command("/bin/sh", "-c", command).CombinedOutput()
	Fatal(fmt.Sprintf("Error returned running a system command: %s", command), err)

	Log(fmt.Sprintf("system.Execute result: %s", string(resultByte)), false)
	return resultByte
}
