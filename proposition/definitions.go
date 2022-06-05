package proposition

import (
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Proposition struct {
	ID        primitive.ObjectID `bson:"_id"`
	Type      string
	DateTime  string
	Desc      string
	Predicate Predicate
	Status    int
	User      int
}

type Predicate struct {
	Label    string
	Value    string
	DataType utils.DataType
	Options  []string
}
