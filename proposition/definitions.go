package proposition

import "github.com/philcantcode/localmapper/utils"

type Proposition struct {
	ID          int
	Type        string
	Date        string
	Description string
	Proposition PropositionItem
	Correction  PropositionItem
	Status      int
	User        int
}

type PropositionItem struct {
	Name     string
	Value    string
	DataType utils.DataType
	Options  []string
}
