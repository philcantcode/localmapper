package utils

type DataType int

const (
	EMPTY DataType = iota
	IP
	IP_RANGE
	MAC
	INTEGER
	DECIMAL
	BOOL
	STRING
	CIDR
	IP6
	MAC6
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
	default:
		ErrorForceFatal("Couldn't do a reverse lookup for DataType (definitions)")
	}

	return "nil"
}
