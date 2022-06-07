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
	Flag     string           // Flag (e.g., -s)
	Value    string           // Set value
	Desc     string           // Contextual info about the flag
	DataType []utils.DataType // e.g., IP, IP Range, String
	Default  string           // Default value that will be used if no value provided
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

func TEST_GENERATE_CAPABILITIES() {

	netBiosScan := Capability{
		Type: "nmap",
		Name: "nbstat NetBIOS",
		Desc: "Attempts to retrieve the target's NetBIOS names and MAC address.",
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Desc:     "UDP Scan",
					Flag:     "-sU",
					DataType: []utils.DataType{utils.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Run Script",
					Flag:     "--script",
					DataType: []utils.DataType{utils.STRING},
					Value:    "nbstat.nse",
					Default:  "nbstat.nse",
				},
				{
					Desc:     "Target",
					Flag:     "",
					DataType: []utils.DataType{utils.CIDR, utils.IP},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Port 137",
					Flag:     "-p137",
					DataType: []utils.DataType{utils.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: []utils.DataType{utils.STRING},
					Value:    "-",
					Default:  "-",
				},
			},
		},
	}

	sysDNSScan := Capability{
		Type: "nmap",
		Name: "System DNS Scan",
		Desc: "Use system DNS resolver configured on this host to identify private hostnames.",
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Desc:     "Disable Port Scan Flag",
					Flag:     "-sn",
					DataType: []utils.DataType{utils.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Target",
					Flag:     "",
					DataType: []utils.DataType{utils.CIDR, utils.IP},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "System DNS Flag",
					Flag:     "--system-dns",
					DataType: []utils.DataType{utils.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: []utils.DataType{utils.STRING},
					Value:    "-",
					Default:  "-",
				},
			},
		},
	}

	pingSweep := Capability{
		Type: "nmap",
		Name: "Ping Sweep",
		Desc: "Perform a discovery Ping Sweep against an IP Range.",
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Desc:     "Disable Port Scan Flag",
					Flag:     "-sn",
					DataType: []utils.DataType{utils.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "IP CIDR Target",
					Flag:     "",
					DataType: []utils.DataType{utils.CIDR},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: []utils.DataType{utils.STRING},
					Value:    "-",
					Default:  "-",
				},
			},
		},
	}

	stealthScan := Capability{
		Type: "nmap",
		Name: "Stealth Scan",
		Desc: "Scan thousands of ports on the target device.",
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Desc:     "Stealth Scan Flag",
					Flag:     "-sS",
					DataType: []utils.DataType{utils.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Disable Ping Flag",
					Flag:     "-Pn",
					DataType: []utils.DataType{utils.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Target",
					Flag:     "",
					DataType: []utils.DataType{utils.CIDR, utils.IP},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: []utils.DataType{utils.STRING},
					Value:    "-",
					Default:  "-",
				},
			},
		},
	}

	osIdent := Capability{
		Type: "nmap",
		Name: "OS Identification Scan",
		Desc: "Attempts to identify the operating system of the host.",
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Desc:     "OS Scan Flag",
					Flag:     "-O",
					DataType: []utils.DataType{utils.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Disable Ping Flag",
					Flag:     "-Pn",
					DataType: []utils.DataType{utils.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Target",
					Flag:     "",
					DataType: []utils.DataType{utils.CIDR, utils.IP},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: []utils.DataType{utils.STRING},
					Value:    "-",
					Default:  "-",
				},
			},
		},
	}

	connectScan := Capability{
		Type: "nmap",
		Name: "TCP Connect Scan",
		Desc: "TCP Connect Scan performs a full connection to the host.",
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Desc:     "Connect Scan Flag",
					Flag:     "-sT",
					DataType: []utils.DataType{utils.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Disable Ping Flag",
					Flag:     "-Pn",
					DataType: []utils.DataType{utils.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "IP Target",
					Flag:     "",
					DataType: []utils.DataType{utils.CIDR, utils.IP},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: []utils.DataType{utils.STRING},
					Value:    "-",
					Default:  "-",
				},
			},
		},
	}

	if utils.UserStringInput("Insert NetBios Scan?") == "y" {
		INSERT_Capability(netBiosScan)
	}

	if utils.UserStringInput("Insert System DNS Scan?") == "y" {
		INSERT_Capability(sysDNSScan)
	}

	if utils.UserStringInput("Insert Ping Sweep Scan?") == "y" {
		INSERT_Capability(pingSweep)
	}

	if utils.UserStringInput("Insert OS Identificaiton Scan?") == "y" {
		INSERT_Capability(osIdent)
	}

	if utils.UserStringInput("Insert Connect Scan?") == "y" {
		INSERT_Capability(connectScan)
	}

	if utils.UserStringInput("Insert Stealth Scan?") == "y" {
		INSERT_Capability(stealthScan)
	}
}
