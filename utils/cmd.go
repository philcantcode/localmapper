package utils

import (
	"os"
	"os/exec"
)

// Run commands from the console return: (output, error)
func Run(command string, args ...string) (string, bool) {
	resultByte, err := exec.Command(command, args...).CombinedOutput()

	if err != nil {
		ErrorHandle("Error returned running a command", err, true)
		return string(resultByte), false
	}

	return string(resultByte), true

}

func RunOnTop(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	ErrorHandle("Couldn't run on top: "+command, err, false)
}
