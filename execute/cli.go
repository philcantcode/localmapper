package execute

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/philcantcode/localmapper/utils"
)

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
