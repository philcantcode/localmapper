package system

import (
	"fmt"
	"os"
	"os/exec"
)

/*
	SetupTools is responsible for extracting required
	resources, setting up tools and checking the various scripts exist.
*/
func SetupTools() {
	if GetConfig("delete-local-resources") == "1" {
		cleanup()
	}

	// TODO: Check that programs exist

	// Acc check, set permissions then execute
	execute("chmod +x ./res/apps/acccheck/acccheck.pl")
	execute("./res/apps/acccheck/acccheck.pl")

	unpackWordlists()
}

/*
	unpackWordlists unzips (7z) the wordlists folder to the
	/localmapper directory.
*/
func unpackWordlists() {
	_, err := os.Stat(GetConfig("wordlist-extract-path"))

	if err == nil {
		return // wordlist already unpacked
	}

	res := execute(fmt.Sprintf("7z x -y %s -o%s", GetConfig("wordlist-path"), GetConfig("external-resources-path")))
	Log(fmt.Sprintf("Unzipped wordlist to %s, result: %s", GetConfig("external-resources-path"), string(res)), true)
}

/*
	execute can perform SUDO privileged commands
*/
func execute(command string) []byte {
	Log(fmt.Sprintf("Attempting system.Execute: %s", command), true)

	resultByte, err := exec.Command("/bin/sh", "-c", command).CombinedOutput()
	Fatal(fmt.Sprintf("Error returned running a system command: %s", command), err)

	Log(fmt.Sprintf("system.Execute result: %s", string(resultByte)), false)
	return resultByte
}

/*
	cleanup removes the /localmapper folder
*/
func cleanup() {
	execute(fmt.Sprintf("rm -rf %s", GetConfig("external-resources-path")))
}
