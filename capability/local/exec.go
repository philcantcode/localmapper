package local

import (
	"fmt"
	"os/exec"

	"github.com/philcantcode/localmapper/system"
)

func Execute(prog string, params []string) []byte {
	system.Log(fmt.Sprintf("Attempting local.Execute: %s > %v", prog, params), true)

	resultByte, err := exec.Command(prog, params...).CombinedOutput()
	system.Fatal(fmt.Sprintf("Error returned in local.Execute running a command: %s > %v", prog, params), err)

	return resultByte
}
