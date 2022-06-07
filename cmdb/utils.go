package cmdb

// FindSysTag returns the found tag or an empty tag
func FindSysTag(label string, entry Entry) (EntryTag, bool, int) {
	for index, entryTag := range entry.SysTags {
		if entryTag.Label == label {
			return entryTag, true, index
		}
	}

	return EntryTag{}, false, -1
}

// FindUsrTag returns the found tag or an empty tag
func FindUsrTag(label string, entry Entry) (EntryTag, bool, int) {
	for index, entryTag := range entry.UsrTags {
		if entryTag.Label == label {
			return entryTag, true, index
		}
	}

	return EntryTag{}, false, -1
}

func RemoveTag(label string, entryTagSet []EntryTag) []EntryTag {
	for index, v := range entryTagSet {
		if v.Label == label {
			entryTagSet = append(entryTagSet[:index], entryTagSet[index+1:]...)
		}
	}

	return entryTagSet
}
