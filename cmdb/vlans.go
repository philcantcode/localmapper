package cmdb

import (
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/utils"
)

func CreateNewVLAN(label, desc, lowIP, highIP string, cmdbType int) {
	if !utils.ValidateIP(lowIP) {
		system.Warning("LowIP not valid in cmdb.HTTP_INSERT_Pending", true)
		return
	}

	if !utils.ValidateIP(highIP) {
		system.Warning("HighIP not valid in cmdb.HTTP_INSERT_Pending", true)
		return
	}

	if !utils.ValidateString(label) {
		system.Warning("Label not valid in cmdb.HTTP_INSERT_Pending", true)
		return
	}

	if !utils.ValidateString(desc) {
		system.Warning("Desc not valid in cmdb.HTTP_INSERT_Pending", true)
		return
	}

	cidrArr, err := utils.IPv4RangeToCIDRRange(lowIP, highIP)
	system.Error("Couldn't create CIDR", err)

	highIpTag := EntityTag{Label: "LowIP", DataType: system.DataType_IP, Values: []string{lowIP}}
	lowIpTag := EntityTag{Label: "HighIP", DataType: system.DataType_IP, Values: []string{highIP}}
	cidr := EntityTag{Label: "CIDR", DataType: system.DataType_CIDR, Values: cidrArr}
	entry := Entity{Label: label, Description: desc, OSILayer: 2, CMDBType: CMDBType(cmdbType), DateSeen: []string{utils.GetDateTime().DateTime}, SysTags: []EntityTag{lowIpTag, highIpTag, cidr}, UsrTags: []EntityTag{}}

	entry.InsertPending()
}
