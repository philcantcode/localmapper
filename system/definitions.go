package system

import "time"

type DataType string
type Category string
type Interpreter string

const (
	DataType_EMPTY         DataType = "EMPTY"
	DataType_IP            DataType = "IP"
	DataType_IP_RANGE      DataType = "IP_RANGE" // 192.168.0.0 - 192.168.0.1
	DataType_MAC           DataType = "MAC"
	DataType_INTEGER       DataType = "INTEGER"
	DataType_DECIMAL       DataType = "DECIMAL"
	DataType_BOOL          DataType = "BOOL" // 1 or 0
	DataType_STRING        DataType = "STRING"
	DataType_CIDR          DataType = "CIDR" // 192.168.0.0/24
	DataType_IP6           DataType = "IP6"
	DataType_MAC6          DataType = "MAC6"
	DataType_CPE           DataType = "CPE"          // Common Platform Enumeration: cpe:/o:linux:linux_kernel:2.6.39
	DataType_CCI           DataType = "CCI"          // Common Capability Identifier
	DataType_CCBI          DataType = "CCBI"         // Common Cookbook Identifier
	DataType_IP_RANGE_LOW  DataType = "IP_RANGE_LOW" // 192.168.0.0 - IP meant to be used in a range calculation but is a single IP
	DataType_IP_RANGE_HIGH DataType = "IP_RANGE_HIGH"
	DataType_USERNAME      DataType = "USERNAME"
	DataType_PORT          DataType = "PORT"     // 1 - 65,535
	DataType_PROTOCOL      DataType = "PROTOCOL" // SSH, SNMP, HTTP, SOCKS5
	DataType_PRODUCT       DataType = "PRODUCT"  // HP Generic Scan Gateway, Amazon Whisperplay DIAL REST servic, Netatalk
	DataType_VENDOR        DataType = "VENDOR"
	DataType_FILE_PATH     DataType = "FILE_PATH"
)

const (
	Category_DDOS       Category = "DDOS"
	Category_DISCOVERY  Category = "DISCOVERY"
	Category_BRUTEFORCE Category = "BRUTEFORCE"
	Category_RESEARCH   Category = "RESEARCH"
)

const (
	Interpreter_UNIVERSAL    Interpreter = "UNIVERSAL"
	Interpreter_NMAP         Interpreter = "NMAP"
	Interpreter_ACCCHECK     Interpreter = "ACCCHECK"
	Interpreter_NBTSCAN      Interpreter = "NBTSCAN"
	Interpreter_SEARCHSPLOIT Interpreter = "SEARCHSPLOIT"
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

type ScheduledJob struct {
	Label     string
	StartTime time.Time
	EndTime   time.Time
	Status    string
}
