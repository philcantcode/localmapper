package capability

import (
	"encoding/json"
	"fmt"

	"github.com/philcantcode/localmapper/capability/nmap"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
)

func ExecuteCapability(capability Capability) []byte {
	system.Log(fmt.Sprintf("Executing Capability: %s", capability.Name), true)

	switch capability.Type {
	case "nmap":
		nmapRun := nmap.Execute(ParamsToArray(capability.Command.Params))
		nmap.INSERT_Nmap(nmapRun)

		result, err := json.Marshal(nmapRun)
		system.Error("Couldn't marshal nmaprun", err)

		return result
	default:
		system.Force(fmt.Sprintf("No capability type to run in Capability.ProcessCapability: %v", capability), true)
		return nil
	}
}

/*
	MatchEntryToCapability determines if a given entry can run a given capability
*/
func MatchEntryToCapability(capability Capability, entry cmdb.Entry) (bool, Capability) {
	var success bool

	for k, capParam := range capability.Command.Params {
		success, capability.Command.Params[k] = MatchParamToTag(capParam, entry.SysTags)

		if !success {
			return false, capability
		}
	}

	return true, capability
}

/*
	MatchParamToTag Determines if given a capability param {"Value": "","DataType": 1, "Default": ""}
	Is there any SysTags that can fulfil the Values
*/
func MatchParamToTag(capParam Param, entryTags []cmdb.EntryTag) (bool, Param) {
	// For each: {DataType.CMDB, DataType.IP}
	for _, pType := range capParam.DataType {
		// If the value is already set, move on
		if capParam.Value != "" {
			return true, capParam
		}

		// Skip empty tags that don't require input
		if pType == system.EMPTY {
			return true, capParam
		}

		for _, eTag := range entryTags {
			// The DataType match
			if pType == eTag.DataType {
				capParam.Value = eTag.Values[len(eTag.Values)-1]
				return true, capParam
			}
		}
	}

	return false, capParam
}

func FirstTimeSetup() {

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
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Run Script",
					Flag:     "--script",
					DataType: []system.DataType{system.STRING},
					Value:    "nbstat.nse",
					Default:  "nbstat.nse",
				},
				{
					Desc:     "Target",
					Flag:     "",
					DataType: []system.DataType{system.CIDR, system.IP},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Port 137",
					Flag:     "-p137",
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: []system.DataType{system.STRING},
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
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Target",
					Flag:     "",
					DataType: []system.DataType{system.CIDR, system.IP},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "System DNS Flag",
					Flag:     "--system-dns",
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: []system.DataType{system.STRING},
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
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "IP CIDR Target",
					Flag:     "",
					DataType: []system.DataType{system.CIDR},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: []system.DataType{system.STRING},
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
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Disable Ping Flag",
					Flag:     "-Pn",
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Target",
					Flag:     "",
					DataType: []system.DataType{system.CIDR, system.IP},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: []system.DataType{system.STRING},
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
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Disable Ping Flag",
					Flag:     "-Pn",
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				}, {
					Desc:     "Slows Down Scan",
					Flag:     "--max-rate",
					DataType: []system.DataType{system.INTEGER},
					Value:    "100",
					Default:  "100",
				},
				{
					Desc:     "Target",
					Flag:     "",
					DataType: []system.DataType{system.CIDR, system.IP},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: []system.DataType{system.STRING},
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
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Disable Ping Flag",
					Flag:     "-Pn",
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "IP Target",
					Flag:     "",
					DataType: []system.DataType{system.CIDR, system.IP},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: []system.DataType{system.STRING},
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
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "ARP Flag",
					Flag:     "-PU",
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "IP Target",
					Flag:     "",
					DataType: []system.DataType{system.CIDR, system.IP},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: []system.DataType{system.STRING},
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
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Stealth Scan Flag",
					Flag:     "-sS",
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Script Flag",
					Flag:     "--script",
					DataType: []system.DataType{system.STRING},
					Value:    "smb-os-discovery.nse",
					Default:  "smb-os-discovery.nse",
				},
				{
					Desc:     "Port Flag",
					Flag:     "-p",
					DataType: []system.DataType{system.STRING},
					Value:    "U:137,T:139",
					Default:  "U:137,T:139",
				},
				{
					Desc:     "IP Target",
					Flag:     "",
					DataType: []system.DataType{system.CIDR, system.IP},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: []system.DataType{system.STRING},
					Value:    "-",
					Default:  "-",
				},
			},
		},
	}

	svcDetection := Capability{
		Type:       "nmap",
		CCI:        "cci:nmap:svc-detection:default",
		Name:       "Service Identification Scan",
		Desc:       "Attempts to identify the service version of running services the host.",
		ResultTags: []string{"OS", "OSGen", "OSVendor", "OSAccuracy"},
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Desc:     "OS Scan Flag",
					Flag:     "-O",
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "Disable Ping Flag",
					Flag:     "-Pn",
					DataType: []system.DataType{system.EMPTY},
					Value:    "",
					Default:  "",
				}, {
					Desc:     "Slows Down Scan",
					Flag:     "--max-rate",
					DataType: []system.DataType{system.INTEGER},
					Value:    "100",
					Default:  "100",
				},
				{
					Desc:     "Target",
					Flag:     "",
					DataType: []system.DataType{system.CIDR, system.IP},
					Value:    "",
					Default:  "",
				},
				{
					Desc:     "XML Output",
					Flag:     "-oX",
					DataType: []system.DataType{system.STRING},
					Value:    "-",
					Default:  "-",
				},
			},
		},
	}

	if len(SELECT_Capability(bson.M{}, bson.M{})) == 0 {
		INSERT_Capability(netBiosScan)
		INSERT_Capability(sysDNSScan)
		INSERT_Capability(pingSweep)
		INSERT_Capability(osIdent)
		INSERT_Capability(connectScan)
		INSERT_Capability(stealthScan)
		INSERT_Capability(arpScan)
		INSERT_Capability(smbScriptScan)
	}
}
