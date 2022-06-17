package capability

import (
	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Capability struct {
	ID            primitive.ObjectID `bson:"_id"`
	CCI           string             // Common Capability Identifier = cci:<tool>:<label>:<edition>
	Command       Command            // Program and params
	Category      system.Category
	Type          string   // e.g., nmap
	Name          string   // Capability name "Ping Sweep"
	Desc          string   // Contextual description
	DisplayFields []string // For hiding results
	ResultTags    []string // The Result tags (e.g., IP, Port) gathered by this capability
	Hidden        bool     // Hidden from the user
}

type Command struct {
	Program string  // e.g., nmap
	Params  []Param // List of parameters to be supplied
}

type Param struct {
	Flag     string            // Flag (e.g., -s)
	Value    string            // Set value
	Desc     string            // Contextual info about the flag
	DataType []system.DataType // e.g., IP, IP Range, String
	Default  string            // Default value that will be used if no value provided
}

func ParamsToArray(params []Param) []string {
	var paramArr []string

	for _, param := range params {
		// If the flag is NOT empty, add the flag
		if param.Flag != "" {
			paramArr = append(paramArr, param.Flag)
		}

		// If the MetaType is NOT 'EMPTY' and the value is NOT empty, add the value
		if param.Value != "" {
			paramArr = append(paramArr, param.Value)
		}
	}

	return paramArr
}
