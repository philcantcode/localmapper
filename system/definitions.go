package system

type DataType int

const (
	EMPTY DataType = iota
	IP
	IP_RANGE // 192.168.0.0 - 192.168.0.1
	MAC
	INTEGER
	DECIMAL
	BOOL // 1 or 0
	STRING
	CIDR // 192.168.0.0/24
	IP6
	MAC6
	CPE          // Common Platform Enumeration: cpe:/o:linux:linux_kernel:2.6.39
	CCI          // Common Capability Identifier
	CCBI         // Common Cookbook Identifier
	IP_RANGE_LOW // 192.168.0.0 - IP meant to be used in a range calculation but is a single IP
	IP_RANGE_HIGH
)

func ReverseDataTypeLookup(datType DataType) string {
	switch datType {
	case 0:
		return "EMPTY"
	case 1:
		return "IP"
	case 2:
		return "IP_RANGE"
	case 3:
		return "MAC"
	case 4:
		return "INTEGER"
	case 5:
		return "DECIMAL"
	case 6:
		return "BOOL"
	case 7:
		return "STRING"
	case 8:
		return "CIDR"
	case 9:
		return "IP6"
	case 10:
		return "MAC6"
	case 11:
		return "CPE"
	default:
		Force("Couldn't do a reverse lookup for DataType (definitions)", true)
	}

	return "nil"
}

func FirstTimeSetup() {
	if len(SELECT_Settings_All()) > 0 {
		return
	}

	settings := []Setting{
		{
			Key:   "version",
			Value: "1.0",
		},
		{
			Key:   "patch",
			Value: "1",
		},
		{
			Key:   "runtime-log",
			Value: "res/logs/runtime.txt",
		},
		{
			Key:   "error-log",
			Value: "res/logs/error.txt",
		},
		{
			Key:   "print-log",
			Value: "res/logs/print.txt",
		},
		{
			Key:   "json-log",
			Value: "res/logs/jsonlog.txt",
		},
		{
			Key:   "db-path",
			Value: "res/database.db",
		},
		{
			Key:   "nmap-path",
			Value: "C:\\Program Files (x86)\\Nmap\\nmap",
		},
		{
			Key:   "server-port",
			Value: "8008",
		},
		{
			Key:   "mongo-enabled",
			Value: "0",
		},
		{
			Key:   "mongo-password-req",
			Value: "0",
		},
		{
			Key:   "mongo-user",
			Value: "root",
		},
		{
			Key:   "mongo-password",
			Value: "rootpassword",
		},
		{
			Key:   "mongo-ip",
			Value: "127.0.0.1",
		},
		{
			Key:   "mongo-port",
			Value: "27017",
		},
		{
			Key:   "date-seen-graph-mins-val",
			Value: "60",
		},
	}

	for _, setting := range settings {
		INSERT_Settings(setting.Key, setting.Value)
	}
}
