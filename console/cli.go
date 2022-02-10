package console

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func RegisterCmdCapability() {
	var command string
	var params []string
	var cmdtype string
	var name string
	var desc string

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter a name for the command (metadata)")
	fmt.Print("[>] ")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("Enter a type for the command (metadata)")
	fmt.Print("[>] ")
	scanner.Scan()
	cmdtype = strings.ToLower(scanner.Text())

	fmt.Println("Enter the base command (e.g., nmap or ping)")
	fmt.Print("[>] ")
	scanner.Scan()
	command = scanner.Text()

	fmt.Println("Enter a series of individual flags (e.g., -sS -Pn -v) one at a time")
	fmt.Println("For user input use notation '<string:ip>' or '<int:port>'")
	fmt.Println("Types Include string:int:port:ip:iprange:email")
	fmt.Println("When finished type 'q'")

	for scanner.Text() != "q" {
		fmt.Print("[>] ")
		scanner.Scan()

		if scanner.Text() != "q" {
			params = append(params, scanner.Text())
		}
	}

	fmt.Println("Enter a description for the command (metadata)")
	fmt.Print("[>] ")
	scanner.Scan()
	desc = scanner.Text()

	var answer string
	fmt.Printf("Is this correct? (y/n): %s %v\n%s\n%s\n%s\n", command, params, name, desc, cmdtype)
	fmt.Print("[>] ")
	scanner.Scan()
	answer = scanner.Text()

	if answer != "y" {
		return
	}

	paramJson, _ := json.Marshal(params)

	id := CheckCommandExists(command, string(paramJson))

	if id == -1 {
		InsertCapability(command, string(paramJson), name, cmdtype, desc)
	}
}

func RunCapability() {
	GetAllCapabilities()
}
