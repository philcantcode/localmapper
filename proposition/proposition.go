package proposition

var propositions = []Proposition{}

func (proposition Proposition) process() {
	if proposition.Type == Proposition_Local_Identity {
		proposition.resolveLocalIPIdentity()
	}

	if proposition.Type == Proposition_IP_Identity_Conflict {
		proposition.ipConflict()
	}
}
