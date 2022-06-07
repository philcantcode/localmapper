package capability

import (
	"fmt"

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
		if param.DataType != utils.EMPTY && param.Value != "" {
			paramArr = append(paramArr, param.Value)
		}

		// If the MetaType and Value are empty, use the default
		if param.DataType != utils.EMPTY && param.Value == "" && param.Default != "" {
			paramArr = append(paramArr, param.Default)
		}

	}

	return paramArr
}

func TEST_NEW_CAPABILITY() {
	cap := Capability{
		Type: "nmap",
		Name: "nbstat NetBIOS",
		Desc: "Attempts to retrieve the target's NetBIOS names and MAC address.",
		Command: Command{
			Program: "nmap",
			Params: []Param{
				Param{
					Desc:     "UDP Scan",
					Flag:     "-sU",
					DataType: utils.EMPTY,
					Value:    "",
					Default:  "",
				},
				Param{
					Desc:     "Run Script",
					Flag:     "--script",
					DataType: utils.STRING,
					Value:    "nbstat.nse",
					Default:  "nbstat.nse",
				},
				Param{
					Desc:     "IP Address Target",
					Flag:     "",
					DataType: utils.IP,
					Value:    "",
					Default:  "",
				},
				Param{
					Desc:     "Port 137",
					Flag:     "-p137",
					DataType: utils.EMPTY,
					Value:    "",
					Default:  "",
				},
				Param{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: utils.STRING,
					Value:    "-",
					Default:  "-",
				},
			},
		},
	}

	fmt.Printf("%+v\n", cap)

	if utils.UserStringInput("Do you want to insert this capability?") == "y" {
		INSERT_Capability(cap)
	}
}
