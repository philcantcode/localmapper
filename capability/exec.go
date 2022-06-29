package capability

import (
	"bufio"
	"os/exec"

	"github.com/philcantcode/localmapper/system"
)

func execute(prog string, params []string) []byte {
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
