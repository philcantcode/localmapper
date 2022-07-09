package cmdb

import (
	"github.com/philcantcode/localmapper/local"
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// FindSysTag returns the found tag or an empty tag
func (entry Entity) FindSysTag(label string) (EntityTag, bool, int) {
	for index, entryTag := range entry.SysTags {
		if entryTag.Label == label {
			return entryTag, true, index
		}
	}

	return EntityTag{}, false, -1
}

// FindUsrTag returns the found tag or an empty tag
func (entry Entity) FindUsrTag(label string) (EntityTag, bool, int) {
	for index, entryTag := range entry.UsrTags {
		if entryTag.Label == label {
			return entryTag, true, index
		}
	}

	return EntityTag{}, false, -1
}

func RemoveTag(label string, entryTagSet []EntityTag) []EntityTag {
	for index, v := range entryTagSet {
		if v.Label == label {
			entryTagSet = append(entryTagSet[:index], entryTagSet[index+1:]...)
		}
	}

	return entryTagSet
}

/*
	CheckSelfIdentity checks to see if the default IP matches the one on file
	if it doesn't, the IP gets updated
*/
func updateSelfIdentity() {
	ip := local.GetDefaultIPGateway().DefaultIP
	local := SELECT_Entities_Joined(bson.M{"systags.label": "Identity"}, bson.M{})

	for _, entity := range local {
		identTag, _, _ := entity.FindSysTag("Identity")
		ipTag, ipTagFound, ipIdx := entity.FindSysTag("IP")
		newEntity := entity

		if utils.ArrayContains("local", identTag.Values) && ipTagFound {

			if !utils.ArrayContains(ip, ipTag.Values) {
				newEntity.SysTags[ipIdx].Values = append(newEntity.SysTags[ipIdx].Values, ip)
			}

			newEntity.SysTags[ipIdx] = newEntity.SysTags[ipIdx].PushToFront(ip)

			newEntity.UPDATE_ENTRY_Inventory()
			newEntity.UPDATE_ENTRY_Pending()
		}
	}
}
