package capability

import (
	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Capability struct {
	ID            primitive.ObjectID `bson:"_id"`
	CCI           string             // Common Capability Identifier = cci:<tool>:<label>:<edition>
	Command       Command            // Program and params
	Type          string             // e.g., nmap
	Name          string             // Capability name "Ping Sweep"
	Desc          string             // Contextual description
	DisplayFields []string           // For hiding results
	ResultTags    []string           // The Result tags (e.g., IP, Port) gathered by this capability
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

func InsertDefaultCapabilities() {

	netBiosScan := Capability{
		Type:       "nmap",
		CCI:        "cci:nmap:nbstat-netbios-script:default",
		Name:       "nbstat NetBIOS",
		Desc:       "Attempts to retrieve the target's NetBIOS names and MAC address.",
		ResultTags: []string{"MAC"},
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
		Type:       "nmap",
		CCI:        "cci:nmap:sys-dns:default",
		Name:       "System DNS Scan",
		Desc:       "Use system DNS resolver configured on this host to identify private hostnames.",
		ResultTags: []string{"HostName"},
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
		Type:       "nmap",
		Name:       "Ping Sweep",
		CCI:        "cci:nmap:ping-sweep:default",
		Desc:       "Perform a discovery Ping Sweep against an IP Range.",
		ResultTags: []string{"IP"},
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
		Type:       "nmap",
		CCI:        "cci:nmap:stealth:default",
		Name:       "Stealth Scan",
		Desc:       "Scan thousands of ports on the target device.",
		ResultTags: []string{"Ports"},
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
		Type:       "nmap",
		CCI:        "cci:nmap:os-ident:default",
		Name:       "OS Identification Scan",
		Desc:       "Attempts to identify the operating system of the host.",
		ResultTags: []string{"OS", "OSGen", "OSVendor", "OSAccuracy"},
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
				}, {
					Desc:     "Slows Down Scan",
					Flag:     "--max-rate",
					DataType: []utils.DataType{utils.INTEGER},
					Value:    "100",
					Default:  "100",
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
		Type:       "nmap",
		CCI:        "cci:nmap:tcp-connect:default",
		Name:       "TCP Connect Scan",
		Desc:       "TCP Connect Scan performs a full connection to the host.",
		ResultTags: []string{"Ports"},
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

	arpScan := Capability{
		Type:       "nmap",
		Name:       "APR Scan",
		CCI:        "cci:nmap:arp:default",
		Desc:       "Perform a scan of the local network using ARP.",
		ResultTags: []string{"IP"},
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
					Desc:     "ARP Flag",
					Flag:     "-PU",
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

	// nmap IP -sU -sS --script smb-os-discovery.nse -p U:137,T:139

	smbScriptScan := Capability{
		Type:       "nmap",
		Name:       "SMB OS Discovery",
		CCI:        "cci:nmap:smb-os-discovery:default",
		Desc:       "Attempts to determine the operating system, computer name, domain, workgroup, and current time over the SMB protocol (ports 445 or 139). This is done by starting a session with the anonymous account, in response to a session starting, the server will send back all this information.",
		ResultTags: []string{"HostName"},
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Desc:     "UDP Scan Flag",
					Flag:     "-sU",
					DataType: []utils.DataType{utils.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Stealth Scan Flag",
					Flag:     "-sS",
					DataType: []utils.DataType{utils.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Script Flag",
					Flag:     "--script",
					DataType: []utils.DataType{utils.STRING},
					Value:    "smb-os-discovery.nse",
					Default:  "smb-os-discovery.nse",
				},
				{
					Desc:     "Port Flag",
					Flag:     "-p",
					DataType: []utils.DataType{utils.STRING},
					Value:    "U:137,T:139",
					Default:  "U:137,T:139",
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

	INSERT_Capability(netBiosScan)
	INSERT_Capability(sysDNSScan)
	INSERT_Capability(pingSweep)
	INSERT_Capability(osIdent)
	INSERT_Capability(connectScan)
	INSERT_Capability(stealthScan)
	INSERT_Capability(arpScan)
	INSERT_Capability(smbScriptScan)

}
