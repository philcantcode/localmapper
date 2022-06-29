package system

import (
	"fmt"
	"os"

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

	setupDirectories()
	unpackWordlists()
	installKatoolin()
	installSearchsploit()

	// execute("katoolin", []string{"2", "1", "22"})

	// Acc check, set permissions
	Execute("chmod +x ./res/apps/acccheck/acccheck.pl")

	// Update tools
	Log("All tools setup correctly", true)
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

	res := Execute(fmt.Sprintf("7z x -y %s -o%s", GetConfig("wordlist-path"), GetConfig("external-resources-path")))
	Log(fmt.Sprintf("Unzipped wordlist to %s, result: %s", GetConfig("external-resources-path"), string(res)), true)
}

func installKatoolin() {
	if utils.DirExists(GetConfig("external-resources-path") + "/katoolin") {
		return
	}

	Execute("add-apt-repository universe -y")
	Execute(fmt.Sprintf("git -C %s clone https://github.com/LionSec/katoolin.git", GetConfig("external-resources-path")))
	Execute(fmt.Sprintf("mv %s/katoolin/katoolin.py /usr/bin/katoolin", GetConfig("external-resources-path")))
	Execute("chmod +x /usr/bin/katoolin")
	Execute("update-alternatives --install /usr/bin/python python /usr/bin/python2 1")
}

func installSearchsploit() {
	installPath := GetConfig("external-resources-path") + "/exploitdb"

	if !utils.DirExists(installPath) {
		Execute(fmt.Sprintf("git -C %s clone https://github.com/offensive-security/exploitdb.git", GetConfig("external-resources-path")))
	}

	Execute(fmt.Sprintf("sed -i 's:/opt/exploitdb:%s:g' %s/.searchsploit_rc", installPath, installPath)) // Replace /opt path with our path
	Execute(fmt.Sprintf("ln -sf %s/searchsploit /usr/bin/searchsploit", installPath))

	go Execute("searchsploit -u")
}

func setupDirectories() {
	//Nmap directory
	path := GetConfig("external-resources-path")

	if !utils.DirExists(path) {
		Execute("mkdir " + path)
	}

	// Nmap directory
	path = GetConfig("nmap-results-dir")

	if !utils.DirExists(path) {
		Execute("mkdir " + path)
	}
}

/*
	cleanup removes the /localmapper folder
*/
func cleanup() {
	Execute(fmt.Sprintf("rm -rf %s", GetConfig("external-resources-path")))
}
