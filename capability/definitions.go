package capability

import (
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Capability struct {
	ID            primitive.ObjectID `bson:"_id"`
	Command       Command            // Program and params
	Type          string             // e.g., nmap
	Name          string             // Capability name "Ping Sweep"
	Desc          string             // Contextual description
	DisplayFields []string           // For hiding results
}

type Command struct {
	Program string  // e.g., nmap
	Params  []Param // List of parameters to be supplied
}

type Param struct {
	Flag     string         // Flag (e.g., -s)
	Value    string         // Set value
	Desc     string         // Contextual info about the flag
	DataType utils.DataType // e.g., IP, IP Range, String
	Default  string         // Default value that will be used if no value provided
}

func ParamsToArray(params []Param) []string {
	var paramArr []string

	for _, param := range params {
		// If the flag is NOT empty, add the flag
		if param.Flag != "" {
			paramArr = append(paramArr, param.Flag)
		}

		// If the MetaType is NOT 'none' and the value is NOT empty, add the value
		if param.DataType != utils.None && param.Value != "" {
			paramArr = append(paramArr, param.Value)
		}

		// If the MetaType and Value are empty, use the default
		if param.DataType != utils.None && param.Value == "" && param.Default != "" {
			paramArr = append(paramArr, param.Default)
		}

	}

	return paramArr
}
