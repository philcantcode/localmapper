package proposition

import (
	"fmt"

	"github.com/philcantcode/localmapper/system"
)

var propositions = []Proposition{}

func (proposition Proposition) resolve() {
	if proposition.Type == Proposition_Local_Identity {
		proposition.resolveLocalIPIdentity()
	}

	if proposition.Type == Proposition_IP_Identity_Conflict {
		proposition.ipConflict()
	}

	for i := range propositions {
		if propositions[i].ID == proposition.ID {
			propositions = append(propositions[:i], propositions[i+1:]...)
			system.Log(fmt.Sprintf("[Proposition Resolved] %s", proposition.Description), true)
		}
	}
}
