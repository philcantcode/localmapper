package ping

import (
	"fmt"

	"github.com/philcantcode/localmapper/adapters/definitions"
	"github.com/philcantcode/localmapper/utils"
)

func RunPingCommand(capability definitions.Capability) {
	utils.Log(fmt.Sprintf("Attempting to Ping: %v", capability), true)

}
