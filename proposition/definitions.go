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
	Evidence    []Evidence
}

type Predicate struct {
	Label    string
	Value    string
	DataType system.DataType
}

type Evidence struct {
	Label string
	Value string
}
