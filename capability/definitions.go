package capability

import (
	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Capability struct {
	ID            primitive.ObjectID `bson:"_id"`
	Label         string             // Capability name "Ping Sweep"
	CCI           string             // Common Capability Identifier = cci:<tool>:<label>:<edition>
	Description   string             // Contextual description
	Category      system.Category
	Interpreter   system.Interpreter
	Hidden        bool     // Hidden from the user
	DisplayFields []string // For hiding results
	ResultTags    []string // The Result tags (e.g., IP, Port) gathered by this capability
	Command       Command  // Program and params
}

type Command struct {
	Program string  // e.g., nmap
	Params  []Param // List of parameters to be supplied
}

type Param struct {
	Flag        string            // Flag (e.g., -s)
	Value       string            // Set value
	Description string            // Contextual info about the flag
	DataType    []system.DataType // e.g., IP, IP Range, String
	Default     string            // Default value that will be used if no value provided
	Options     []string          // Options the user may choose from
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
