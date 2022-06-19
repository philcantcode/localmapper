package proposition

import (
	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Proposition struct {
	ID          primitive.ObjectID `bson:"_id"`
	Type        string
	DateTime    string
	Description string
	Predicate   Predicate
	Status      int // 0 = Open, 1 = Accepted, 2 = Deleted
	User        int
}

type Predicate struct {
	Label    string
	Value    string
	DataType system.DataType
	Options  []string
}
