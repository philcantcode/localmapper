package ping

import (
	"fmt"

	"github.com/philcantcode/localmapper/capability"
	"github.com/philcantcode/localmapper/utils"
)

func RunPingCommand(capability capability.Capability) {
	utils.Log(fmt.Sprintf("Attempting to Ping: %v", capability), true)

}
