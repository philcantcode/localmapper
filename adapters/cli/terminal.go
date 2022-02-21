package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/philcantcode/localmapper/application/localhost"
	"github.com/philcantcode/localmapper/utils"
)

func InitCLI() {
	switch utils.Configs["MODE"] {
	case "Interactive":
		go interactiveCLI()
	}
}

func interactiveCLI() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Text() != "q" {
		fmt.Print("[>] ")
		scanner.Scan()
		runCMD(scanner.Text())
	}
}

func runCMD(cmd string) {
	switch cmd {
	case "ip":
		utils.PrettyPrint(localhost.IpInfo())
	case "os":
		utils.PrettyPrint(localhost.OSInfo())
	case "run capability":
		runCapability()
	case "help":
		fmt.Println("Available Commands: ip, os, run capability, help")
	}
}
