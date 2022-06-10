package ping

import (
	"fmt"

	"github.com/philcantcode/localmapper/capability"
	"github.com/philcantcode/localmapper/system"
)

func RunPingCommand(capability capability.Capability) {
	system.Log(fmt.Sprintf("Attempting to Ping: %v", capability), true)

}
