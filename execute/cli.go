package execute

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/philcantcode/localmapper/utils"
)

func RegisterCmdCapability() {
	var params []string
	var cmdtype string
	var name string
	var desc string
	var interpreter string

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
	params = append(params, scanner.Text())

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

	fmt.Println("Enter an interpreter string for the command (e.g., string:json)")
	fmt.Print("[>] ")
	scanner.Scan()
	interpreter = strings.ToLower(scanner.Text())

	var answer string
	fmt.Printf("Is this correct? (y/n): %v\n%s\n%s\n%s\n%s\n", params, name, desc, cmdtype, interpreter)
	fmt.Print("[>] ")
	scanner.Scan()
	answer = scanner.Text()

	if answer != "y" {
		return
	}

	paramJson, _ := json.Marshal(params)

	id := CheckCommandExists(string(paramJson))

	if id == -1 {
		InsertCapability(string(paramJson), name, cmdtype, desc, interpreter)
	}
}

func RunCapability() {
	capabilities := GetAllCapabilities()
	scanner := bufio.NewScanner(os.Stdin)
	var capID int

	fmt.Println("Enter a capability number: ")

	for _, k := range capabilities {
		fmt.Printf("    %d - %s, %s (%s)\n", k.ID, k.Name, k.Desc, k.Type)
	}

	fmt.Print("[>] ")
	scanner.Scan()
	capID, _ = strconv.Atoi(scanner.Text())

	for _, k := range capabilities {
		if capID == k.ID {
			result, success := Run(k.Interpreter, k.Params[0], swapOutCapabilityParamsWithCLIValues(k.Params[1:])...)

			utils.Log(fmt.Sprintf("Command status: %v\n", success), true)
			utils.PrettyPrint(result)
		}
	}

	utils.Log("RunCapability done", true)
}

func swapOutCapabilityParamsWithCLIValues(params []string) []string {
	scanner := bufio.NewScanner(os.Stdin)

	for i, param := range params {
		param = strings.Replace(param, "<", "", -1)
		param = strings.Replace(param, ">", "", -1)

		switch param {
		case "string":
			fmt.Printf("Please enter a String for the parameter: %s\n", param)
			scanner.Scan()
			params[i] = scanner.Text()
		case "string:iprange":
			fmt.Printf("Please enter an IP Range for the parameter: %s\n", param)
			scanner.Scan()
			params[i] = scanner.Text()
		case "string:ip":
			fmt.Printf("Please enter an IP Address string for the parameter: %s\n", param)
			scanner.Scan()
			params[i] = scanner.Text()
		}
	}

	return params
}
