package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/philcantcode/localmapper/application/database"
	"github.com/philcantcode/localmapper/application/nmap"
	"github.com/philcantcode/localmapper/utils"
)

func runCapability() {
	capabilities := database.SelectAllCapabilities()
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
			capability := k
			k.Params = swapOutCapabilityParamsWithCLIValues(k.Params[1:])

			switch capability.Type {
			case "nmap":
				nmapRun := nmap.RunNmapCommand(capability)
				utils.PrettyPrint(nmapRun)
			default:
				utils.ErrorForceFatal("No capability type to run in RunCapability")
			}

			utils.Log(fmt.Sprintf("Capability Complete: [%s] %s", k.Type, k.Name), true)
			return
		}
	}

	utils.ErrorForceFatal("Could not find a patching capability")
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
