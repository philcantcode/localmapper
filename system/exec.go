package system

import (
	"bufio"
	"fmt"
	"os/exec"
)

func Execute(command string) []byte {
	Log("system.execute running: "+command, true)

	cmd := exec.Command("/bin/sh", "-c", command)
	cmdReader, err := cmd.StdoutPipe()
	Fatal("Couldn't prepare exec.Command", err)

	resultBytes := []byte{}
	scanner := bufio.NewScanner(cmdReader)

	go func() {
		for scanner.Scan() {
			resByte := scanner.Bytes()

			resultBytes = append(resultBytes, resByte...)
			fmt.Println(string(resByte))
		}
	}()

	err = cmd.Start()
	Fatal("Couldn't start exec.Start", err)

	err = cmd.Wait()
	Error("Couldn't wait exec.Wait", err)

	return resultBytes
}
