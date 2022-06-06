package cmdb

// FindSysTag returns the found tag or an empty tag
func FindSysTag(label string, entry Entry) (EntryTag, bool) {
	for _, entryTag := range entry.SysTags {
		if entryTag.Label == label {
			return entryTag, true
		}
	}

	return EntryTag{}, false
}

// FindUsrTag returns the found tag or an empty tag
func FindUsrTag(label string, entry Entry) (EntryTag, bool) {
	for _, entryTag := range entry.UsrTags {
		if entryTag.Label == label {
			return entryTag, true
		}
	}

	return EntryTag{}, false
}

func RemoveTag(entryTagSet []EntryTag, label string) []EntryTag {
	var result []EntryTag

	for index, v := range entryTagSet {
		if v.Label == label {
			result = append(entryTagSet[:index], entryTagSet[index+1:]...)
		}
	}

	return result
}
