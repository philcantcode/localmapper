package ping

import (
	"fmt"

	"github.com/philcantcode/localmapper/adapters/blueprint"
	"github.com/philcantcode/localmapper/utils"
)

func RunPingCommand(capability blueprint.Capability) {
	utils.Log(fmt.Sprintf("Attempting to Ping: %v", capability), true)

}
