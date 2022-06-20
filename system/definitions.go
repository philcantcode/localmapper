package system

type DataType string
type Category string
type Interpreter string

const (
	EMPTY         DataType = "EMPTY"
	IP            DataType = "IP"
	IP_RANGE      DataType = "IP_RANGE" // 192.168.0.0 - 192.168.0.1
	MAC           DataType = "MAC"
	INTEGER       DataType = "INTEGER"
	DECIMAL       DataType = "DECIMAL"
	BOOL          DataType = "BOOL" // 1 or 0
	STRING        DataType = "STRING"
	CIDR          DataType = "CIDR" // 192.168.0.0/24
	IP6           DataType = "IP6"
	MAC6          DataType = "MAC6"
	CPE           DataType = "CPE"          // Common Platform Enumeration: cpe:/o:linux:linux_kernel:2.6.39
	CCI           DataType = "CCI"          // Common Capability Identifier
	CCBI          DataType = "CCBI"         // Common Cookbook Identifier
	IP_RANGE_LOW  DataType = "IP_RANGE_LOW" // 192.168.0.0 - IP meant to be used in a range calculation but is a single IP
	IP_RANGE_HIGH DataType = "IP_RANGE_HIGH"
	USERNAME      DataType = "USERNAME"
)

const (
	DDOS      Category = "DDOS"
	DISCOVERY Category = "DISCOVERY"
)

const (
	UNIVERSAL Interpreter = "UNIVERSAL"
	NMAP      Interpreter = "NMAP"
	ACCCHECK  Interpreter = "ACCCHECK"
	NBTSCAN   Interpreter = "NBTSCAN"
)

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
