package capability

import (
	"encoding/json"
	"fmt"

	"github.com/philcantcode/localmapper/capability/nmap"
	"github.com/philcantcode/localmapper/utils"
)

func ProcessCapability(capability Capability) []byte {
	switch capability.Type {
	case "nmap":
		nmapRun := nmap.Execute(ParamsToArray(capability.Command.Params))
		nmap.InsertNetworkNmap(nmapRun)
		utils.PrintLog(utils.PrettyPrintToStr(nmapRun))

		result, err := json.Marshal(nmapRun)
		utils.ErrorLog("Couldn't marshal nmaprun", err, true)

		return result
	default:
		utils.ErrorForceFatal(fmt.Sprintf("No capability type to run in Capability.ProcessCapability: %v", capability))
		return nil
	}
}
