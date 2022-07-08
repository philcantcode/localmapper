package proposition

import (
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func Init() {
	setupLocalIPIdentity()
	recalcualteVlanCIDR()

}

func recalcualteVlanCIDR() {
	entries := cmdb.SELECT_ENTRY_Inventory(bson.M{"cmdbtype": cmdb.VLAN}, bson.M{})

	for _, entry := range entries {
		// Check CMDB entry is of type VLAN
		if entry.CMDBType != cmdb.VLAN {
			continue
		}

		lowIP, lowFound, _ := entry.FindSysTag("LowIP")
		highIP, highFound, _ := entry.FindSysTag("HighIP")

		// Check that both of the user tags for the IPs are set
		if !lowFound && !highFound {
			continue
		}

		cidr, err := utils.IPv4RangeToCIDRRange(lowIP.Values[0], highIP.Values[0])
		system.Error("Couldn't generate CIDR for: "+entry.Label, err)

		// Remove old CMDB tags so new one can be calcualted
		_, found, index := entry.FindSysTag("CIDR")

		if found {
			entry.SysTags[index] = cmdb.EntityTag{Label: "CIDR", Description: "CIDR range for this VLAN.", DataType: system.DataType_CIDR, Values: cidr}
		} else {
			entry.SysTags = append(entry.SysTags, cmdb.EntityTag{Label: "CIDR", Description: "CIDR range for this VLAN.", DataType: system.DataType_CIDR, Values: cidr})
		}

		entry.UPDATE_ENTRY_Inventory()
	}
}
