package local

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/philcantcode/localmapper/system"
)

// func Execute(prog string, params []string) []byte {
// 	system.Log(fmt.Sprintf("Attempting local.Execute: %s > %v", prog, params), true)

// 	resultByte, err := exec.Command(prog, params...).CombinedOutput()
// 	system.Fatal(fmt.Sprintf("Error returned in local.Execute running a command: %s > %v", prog, params), err)

// 	return resultByte
// }

func Execute(prog string, params []string) []byte {
	system.Log(fmt.Sprintf("Attempting local.Execute: %s > %v", prog, params), true)

	cmd := exec.Command(prog, params...)
	cmdReader, err := cmd.StdoutPipe()
	resultBytes := []byte{}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(0)
	}

	scanner := bufio.NewScanner(cmdReader)

	go func() {
		for scanner.Scan() {
			resByte := scanner.Bytes()

			resultBytes = append(resultBytes, resByte...)
			//fmt.Printf("\t > %s\n", string(resByte))
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(0)
	}

	err = cmd.Wait()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting Cmd", err)
		os.Exit(0)
	}

	//system.Fatal(fmt.Sprintf("Error returned in local.Execute running a command: %s > %v", prog, params), err)

	return resultBytes
}
