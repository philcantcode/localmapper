package capability

import (
	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
)

func FirstTimeSetup() {

	// smurf6 := Capability{
	// 	Type:       "thc-ipv6",
	// 	CCI:        "cci:thc-ipv6:smurf6:default",
	// 	Name:       "Smurf6 Denial of Service (DoS)",
	// 	Desc:       "Attempts to DoS a host.",
	// 	Category:   system.DDOS,
	// 	ResultTags: []string{},
	// 	Command: Command{
	// 		Program: "atk6-smurf6",
	// 		Params: []Param{
	// 			{
	// 				Desc:     "UDP Scan",
	// 				Flag:     "-sU",
	// 				DataType: []system.DataType{system.EMPTY},
	// 				Value:    "",
	// 				Default:  "",
	// 			},
	// 			{
	// 				Desc:     "Run Script",
	// 				Flag:     "--script",
	// 				DataType: []system.DataType{system.STRING},
	// 				Value:    "nbstat.nse",
	// 				Default:  "nbstat.nse",
	// 			},
	// 			{
	// 				Desc:     "Target",
	// 				Flag:     "",
	// 				DataType: []system.DataType{system.CIDR, system.IP},
	// 				Value:    "",
	// 				Default:  "",
	// 			},
	// 			{
	// 				Desc:     "Port 137",
	// 				Flag:     "-p137",
	// 				DataType: []system.DataType{system.EMPTY},
	// 				Value:    "",
	// 				Default:  "",
	// 			},
	// 			{
	// 				Desc:     "XML Output",
	// 				Flag:     "-oX",
	// 				DataType: []system.DataType{system.STRING},
	// 				Value:    "-",
	// 				Default:  "-",
	// 			},
	// 		},
	// 	},
	// }

	// accCheck := Capability{
	// 	Type:       "acccheck",
	// 	CCI:        "cci:kali:acccheck:default",
	// 	Name:       "AC-Check Auth Guessing",
	// 	Desc:       "Connects to IPC and ADMIN shares using SMB to attempt logins.",
	// 	Category:   system.DDOS,
	// 	ResultTags: []string{},
	// 	Command: Command{
	// 		Program: "acccheck.pl",
	// 		Params: []Param{
	// 			{
	// 				Desc:     "Target",
	// 				Flag:     "-t",
	// 				DataType: []system.DataType{system.IP},
	// 				Value:    "",
	// 				Default:  "",
	// 			},
	// 			{
	// 				Desc:     "Run Script",
	// 				Flag:     "--script",
	// 				DataType: []system.DataType{system.STRING},
	// 				Value:    "nbstat.nse",
	// 				Default:  "nbstat.nse",
	// 			},
	// 			{
	// 				Desc:     "Target",
	// 				Flag:     "",
	// 				DataType: []system.DataType{system.CIDR, system.IP},
	// 				Value:    "",
	// 				Default:  "",
	// 			},
	// 			{
	// 				Desc:     "Port 137",
	// 				Flag:     "-p137",
	// 				DataType: []system.DataType{system.EMPTY},
	// 				Value:    "",
	// 				Default:  "",
	// 			},
	// 			{
	// 				Desc:     "XML Output",
	// 				Flag:     "-oX",
	// 				DataType: []system.DataType{system.STRING},
	// 				Value:    "-",
	// 				Default:  "-",
	// 			},
	// 		},
	// 	},
	// }

	netBiosScan := Capability{
		Type:        "nmap",
		CCI:         "cci:nmap:nbstat-netbios-script:default",
		Label:       "nbstat NetBIOS",
		Description: "Attempts to retrieve the target's NetBIOS names and MAC address.",
		ResultTags:  []string{"MAC"},
		Category:    system.DISCOVERY,
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Description: "UDP Scan",
					Flag:        "-sU",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				},
				{
					Description: "Run Script",
					Flag:        "--script",
					DataType:    []system.DataType{system.STRING},
					Value:       "nbstat.nse",
					Default:     "nbstat.nse",
				},
				{
					Description: "Target",
					Flag:        "",
					DataType:    []system.DataType{system.CIDR, system.IP},
					Value:       "",
					Default:     "",
				},
				{
					Description: "Port 137",
					Flag:        "-p137",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				},
				{
					Description: "XML Output",
					Flag:        "-oX",
					DataType:    []system.DataType{system.STRING},
					Value:       "-",
					Default:     "-",
				},
			},
		},
	}

	sysDNSScan := Capability{
		Type:        "nmap",
		CCI:         "cci:nmap:sys-dns:default",
		Label:       "System DNS Scan",
		Description: "Use system DNS resolver configured on this host to identify private hostnames.",
		ResultTags:  []string{"HostName"},
		Category:    system.DISCOVERY,
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Description: "Disable Port Scan Flag",
					Flag:        "-sn",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				},
				{
					Description: "Target",
					Flag:        "",
					DataType:    []system.DataType{system.CIDR, system.IP},
					Value:       "",
					Default:     "",
				},
				{
					Description: "System DNS Flag",
					Flag:        "--system-dns",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				},
				{
					Description: "XML Output",
					Flag:        "-oX",
					DataType:    []system.DataType{system.STRING},
					Value:       "-",
					Default:     "-",
				},
			},
		},
	}

	pingSweep := Capability{
		Type:        "nmap",
		Label:       "Ping Sweep",
		CCI:         "cci:nmap:ping-sweep:default",
		Description: "Perform a discovery Ping Sweep against an IP Range.",
		ResultTags:  []string{"IP"},
		Category:    system.DISCOVERY,
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Description: "Disable Port Scan Flag",
					Flag:        "-sn",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				},
				{
					Description: "IP CIDR Target",
					Flag:        "",
					DataType:    []system.DataType{system.CIDR},
					Value:       "",
					Default:     "",
				},
				{
					Description: "XML Output",
					Flag:        "-oX",
					DataType:    []system.DataType{system.STRING},
					Value:       "-",
					Default:     "-",
				},
			},
		},
	}

	stealthScan := Capability{
		Type:        "nmap",
		CCI:         "cci:nmap:stealth:default",
		Label:       "Stealth Scan",
		Description: "Scan thousands of ports on the target device.",
		ResultTags:  []string{"Ports"},
		Category:    system.DISCOVERY,
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Description: "Stealth Scan Flag",
					Flag:        "-sS",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				},
				{
					Description: "Disable Ping Flag",
					Flag:        "-Pn",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				},
				{
					Description: "Target",
					Flag:        "",
					DataType:    []system.DataType{system.CIDR, system.IP},
					Value:       "",
					Default:     "",
				},
				{
					Description: "XML Output",
					Flag:        "-oX",
					DataType:    []system.DataType{system.STRING},
					Value:       "-",
					Default:     "-",
				},
			},
		},
	}

	osIdent := Capability{
		Type:        "nmap",
		CCI:         "cci:nmap:os-ident:default",
		Label:       "OS Identification Scan",
		Description: "Attempts to identify the operating system of the host.",
		ResultTags:  []string{"OS", "OSGen", "OSVendor", "OSAccuracy"},
		Category:    system.DISCOVERY,
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Description: "OS Scan Flag",
					Flag:        "-O",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				},
				{
					Description: "Disable Ping Flag",
					Flag:        "-Pn",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				}, {
					Description: "Slows Down Scan",
					Flag:        "--max-rate",
					DataType:    []system.DataType{system.INTEGER},
					Value:       "100",
					Default:     "100",
				},
				{
					Description: "Target",
					Flag:        "",
					DataType:    []system.DataType{system.CIDR, system.IP},
					Value:       "",
					Default:     "",
				},
				{
					Description: "XML Output",
					Flag:        "-oX",
					DataType:    []system.DataType{system.STRING},
					Value:       "-",
					Default:     "-",
				},
			},
		},
	}

	connectScan := Capability{
		Type:        "nmap",
		CCI:         "cci:nmap:tcp-connect:default",
		Label:       "TCP Connect Scan",
		Description: "TCP Connect Scan performs a full connection to the host.",
		ResultTags:  []string{"Ports"},
		Category:    system.DISCOVERY,
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Description: "Connect Scan Flag",
					Flag:        "-sT",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				},
				{
					Description: "Disable Ping Flag",
					Flag:        "-Pn",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				},
				{
					Description: "IP Target",
					Flag:        "",
					DataType:    []system.DataType{system.CIDR, system.IP},
					Value:       "",
					Default:     "",
				},
				{
					Description: "XML Output",
					Flag:        "-oX",
					DataType:    []system.DataType{system.STRING},
					Value:       "-",
					Default:     "-",
				},
			},
		},
	}

	arpScan := Capability{
		Type:        "nmap",
		Label:       "APR Scan",
		CCI:         "cci:nmap:arp:default",
		Description: "Perform a scan of the local network using ARP.",
		ResultTags:  []string{"IP"},
		Category:    system.DISCOVERY,
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Description: "Disable Port Scan Flag",
					Flag:        "-sn",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				},
				{
					Description: "ARP Flag",
					Flag:        "-PU",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				},
				{
					Description: "IP Target",
					Flag:        "",
					DataType:    []system.DataType{system.CIDR, system.IP},
					Value:       "",
					Default:     "",
				},
				{
					Description: "XML Output",
					Flag:        "-oX",
					DataType:    []system.DataType{system.STRING},
					Value:       "-",
					Default:     "-",
				},
			},
		},
	}

	// nmap IP -sU -sS --script smb-os-discovery.nse -p U:137,T:139

	smbScriptScan := Capability{
		Type:        "nmap",
		Label:       "SMB OS Discovery",
		CCI:         "cci:nmap:smb-os-discovery:default",
		Description: "Attempts to determine the operating system, computer name, domain, workgroup, and current time over the SMB protocol (ports 445 or 139). This is done by starting a session with the anonymous account, in response to a session starting, the server will send back all this information.",
		ResultTags:  []string{"HostName"},
		Category:    system.DISCOVERY,
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Description: "UDP Scan Flag",
					Flag:        "-sU",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				},
				{
					Description: "Stealth Scan Flag",
					Flag:        "-sS",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				},
				{
					Description: "Script Flag",
					Flag:        "--script",
					DataType:    []system.DataType{system.STRING},
					Value:       "smb-os-discovery.nse",
					Default:     "smb-os-discovery.nse",
				},
				{
					Description: "Port Flag",
					Flag:        "-p",
					DataType:    []system.DataType{system.STRING},
					Value:       "U:137,T:139",
					Default:     "U:137,T:139",
				},
				{
					Description: "IP Target",
					Flag:        "",
					DataType:    []system.DataType{system.CIDR, system.IP},
					Value:       "",
					Default:     "",
				},
				{
					Description: "XML Output",
					Flag:        "-oX",
					DataType:    []system.DataType{system.STRING},
					Value:       "-",
					Default:     "-",
				},
			},
		},
	}

	svcDetection := Capability{
		Type:        "nmap",
		CCI:         "cci:nmap:svc-detection:default",
		Label:       "Service Identification Scan",
		Description: "Attempts to identify the service version of running services the host.",
		ResultTags:  []string{"OS", "OSGen", "OSVendor", "OSAccuracy"},
		Category:    system.DISCOVERY,
		Command: Command{
			Program: "nmap",
			Params: []Param{
				{
					Description: "OS Scan Flag",
					Flag:        "-sV",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				},
				{
					Description: "Disable Ping Flag",
					Flag:        "-Pn",
					DataType:    []system.DataType{system.EMPTY},
					Value:       "",
					Default:     "",
				}, {
					Description: "Slows Down Scan",
					Flag:        "--max-rate",
					DataType:    []system.DataType{system.INTEGER},
					Value:       "100",
					Default:     "100",
				},
				{
					Description: "Target",
					Flag:        "",
					DataType:    []system.DataType{system.CIDR, system.IP},
					Value:       "",
					Default:     "",
				},
				{
					Description: "XML Output",
					Flag:        "-oX",
					DataType:    []system.DataType{system.STRING},
					Value:       "-",
					Default:     "-",
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
		INSERT_Capability(svcDetection)
	}
}
