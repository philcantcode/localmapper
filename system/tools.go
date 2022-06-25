package system

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/philcantcode/localmapper/utils"
)

var deleteLocalResources = false

/*
	SetupTools is responsible for extracting required
	resources, setting up tools and checking the various scripts exist.
*/
func SetupTools() {
	if GetConfig("delete-local-resources") == "1" || deleteLocalResources {
		cleanup()
	}

	// TODO: Check that programs exist

	unpackWordlists()
	installKatoolin()
	installSearchsploit()
	setupNmapDir()

	// execute("katoolin", []string{"2", "1", "22"})

	// Acc check, set permissions then execute
	execute("chmod +x ./res/apps/acccheck/acccheck.pl", nil)
	execute("./res/apps/acccheck/acccheck.pl", nil)

	// Update tools
	Log("All tools setup correctly, running updates concurrently", true)
	go execute("searchsploit -u", nil)
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

	res := execute(fmt.Sprintf("7z x -y %s -o%s", GetConfig("wordlist-path"), GetConfig("external-resources-path")), nil)
	Log(fmt.Sprintf("Unzipped wordlist to %s, result: %s", GetConfig("external-resources-path"), string(res)), true)
}

func installKatoolin() {
	if utils.DirExists(GetConfig("external-resources-path") + "/katoolin") {
		return
	}

	execute("add-apt-repository universe -y", nil)
	execute(fmt.Sprintf("git -C %s clone https://github.com/LionSec/katoolin.git", GetConfig("external-resources-path")), nil)
	execute(fmt.Sprintf("mv %s/katoolin/katoolin.py /usr/bin/katoolin", GetConfig("external-resources-path")), nil)
	execute("chmod +x /usr/bin/katoolin", nil)
	execute("update-alternatives --install /usr/bin/python python /usr/bin/python2 1", nil)
}

func installSearchsploit() {
	installPath := GetConfig("external-resources-path") + "/exploitdb"

	if !utils.DirExists(installPath) {
		execute(fmt.Sprintf("git -C %s clone https://github.com/offensive-security/exploitdb.git", GetConfig("external-resources-path")), nil)
	}

	execute(fmt.Sprintf("sed -i 's:/opt/exploitdb:%s:g' %s/.searchsploit_rc", installPath, installPath), nil) // Replace /opt path with our path
	execute(fmt.Sprintf("ln -sf %s/searchsploit /usr/bin/searchsploit", installPath), nil)
}

func setupNmapDir() {
	path := GetConfig("nmap-results-dir")

	if utils.DirExists(path) {
		return
	}

	execute("mkdir "+path, nil)
}

/*
	cleanup removes the /localmapper folder
*/
func cleanup() {
	execute(fmt.Sprintf("rm -rf %s", GetConfig("external-resources-path")), nil)
}

func execute(command string, ctrlParams []string) []byte {
	Log("system.execute running: "+command, true)

	cmd := exec.Command("/bin/sh", "-c", command)
	cmdReader, err := cmd.StdoutPipe()
	cmdWriter, err := cmd.StdinPipe()

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
			fmt.Println(string(resByte))
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(0)
	}

	if ctrlParams != nil {
		for _, param := range ctrlParams {
			Log("Running parameter: "+param, true)
			cmdWriter.Write([]byte(param))
			time.Sleep(2 * time.Second)
		}

		cmd.Process.Signal(os.Kill)
	}

	err = cmd.Wait()

	//system.Fatal(fmt.Sprintf("Error returned in local.Execute running a command: %s > %v", prog, params), err)

	return resultBytes
}
