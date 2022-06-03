package ping

import (
	"fmt"

	"github.com/philcantcode/localmapper/core"
	"github.com/philcantcode/localmapper/utils"
)

func RunPingCommand(capability core.Capability) {
	utils.Log(fmt.Sprintf("Attempting to Ping: %v", capability), true)

}
