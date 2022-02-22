package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/philcantcode/localmapper/adapters/blueprint"
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
			capability = setCapabilityParamValues(capability)

			switch capability.Type {
			case "nmap":
				nmapRun := nmap.RunNmapCommand(capability)
				utils.PrettyPrint(nmapRun)
			default:
				utils.ErrorForceFatal("No capability type to run in RunCapability")
			}

			utils.Log(fmt.Sprintf("Capability Complete: [%s] %s", capability.Type, capability.Name), true)
			return
		}
	}

	utils.ErrorForceFatal("Could not find a patching capability")
}

func setCapabilityParamValues(capability blueprint.Capability) blueprint.Capability {
	scanner := bufio.NewScanner(os.Stdin)

	capabilityCopy := capability

	for i, param := range capability.Command.Params {

		fmt.Println(param)
		switch param.MetaType {
		//TODO: Add input type validation for blueprint.DataType
		case blueprint.None:
			continue
		default:
			fmt.Printf("Please input a type of: %s\n", blueprint.ReverseDataTypeLookup(param.MetaType))
			scanner.Scan()
			capabilityCopy.Command.Params[i].Value = scanner.Text()
		}
	}

	return capabilityCopy
}
