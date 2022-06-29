package local

import (
	"bufio"
	"fmt"
	"os/exec"
	"time"

	"github.com/philcantcode/localmapper/system"
)

type ExecResult struct {
	Label    string
	DateTime time.Time
	Output   []string
	Result   int
}

func Execute(prog string, params []string) []byte {
	system.Log(fmt.Sprintf("Attempting local.Execute: %s %v", prog, params), true)

	cmd := exec.Command(prog, params...)
	cmdReader, err := cmd.StdoutPipe()
	system.Fatal("Couldn't prepare command", err)

	resultBytes := []byte{}

	scanner := bufio.NewScanner(cmdReader)

	go func() {
		for scanner.Scan() {
			resByte := scanner.Bytes()

			resultBytes = append(resultBytes, resByte...)
			//fmt.Printf("\t > %s\n", string(resByte))
		}
	}()

	err = cmd.Start()
	system.Fatal("Couldn't start command", err)

	err = cmd.Wait()
	system.Fatal("Couldn't wait on command", err)

	return resultBytes
}
