package nbtscan

import (
	"context"
	"fmt"

	"github.com/philcantcode/localmapper/system"
)

func (nbtscan NBTScan) Insert() {
	system.Log("Attempting to INSERT_NbtScan", false)

	insertResult, err := system.Results_Nbscan_DB.InsertOne(context.Background(), nbtscan)

	system.Fatal("Couldn't INSERT_NbtScan", err)
	system.Log(fmt.Sprintf("New Insert at: %s", insertResult), true)
}
