package proposition

import (
	"github.com/philcantcode/localmapper/system"
)

type PropType string

const (
	Proposition_Local_Identity       PropType = "LOCAL_IDENTITY"
	Proposition_IP_Identity_Conflict PropType = "IP_IDENTITY_CONFLICT"
)

type Proposition struct {
	ID          string
	Type        PropType
	DateTime    string
	Description string
	Predicates  []Predicate
}

type Predicate struct {
	Label    string
	Value    string
	DataType system.DataType
	Chosen   bool
}

func (prop Proposition) getChosen() Predicate {
	for _, pred := range prop.Predicates {
		if pred.Chosen {
			return pred
		}
	}

	system.Warning("No predicate chosen", true)
	return Predicate{}
}
