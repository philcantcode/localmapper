package proposition

import (
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/utils"
)

var Propositions = []Proposition{}

func (prop Proposition) Push() {
	prop.ID = utils.HashStruct(prop)

	// Remove duplicates by ID
	for _, existing := range Propositions {
		if prop.ID == existing.ID {
			system.Log("Forgoing adding proposition because matching ID exists", false)
			return
		}
	}

	// Remove duplicates by evidence
	for _, existing := range Propositions {
		if len(existing.Evidence) == len(prop.Evidence) {
			for i, existEvid := range existing.Evidence {
				if existEvid == prop.Evidence[i] {
					system.Log("Forgoing adding proposition because matching evidence exists", false)
					return
				}
			}
		}
	}

	Propositions = append(Propositions, prop)
}

/*
	Checks whether a given PropType exists in the list
	of propositions, useful for props where there should
	only be one of them at a time.
*/
func CheckPropTypeExists(ptype PropType) bool {
	for _, prop := range Propositions {
		if prop.Type == Proposition_Local_Identity {
			return true
		}
	}

	return false
}

func Pop(idx int) {
	Propositions = append(Propositions[:idx], Propositions[idx+1:]...)
}

func (prop Proposition) GetEvidenceValue(label string) string {
	for _, evidence := range prop.Evidence {
		if evidence.Label == label {
			return evidence.Value
		}
	}

	system.Warning("Couldn't find a matching proposition evidence value from label", true)
	return ""
}
