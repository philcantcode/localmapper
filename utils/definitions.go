package utils

type DataType int

const (
	String DataType = iota
	IP
	IPRange
	MAC
	Integer
	Decimal
	Bool
	None
	CIDR
)

func ReverseDataTypeLookup(datType DataType) string {
	switch datType {
	case 0:
		return "String"
	case 1:
		return "IP"
	case 2:
		return "IPRange"
	case 3:
		return "MAC"
	case 4:
		return "Integer"
	case 5:
		return "Decimal"
	case 6:
		return "Bool"
	case 7:
		return "None"
	case 9:
		return "CIDR"
	default:
		ErrorForceFatal("Couldn't do a reverse lookup for DataType (definitions)")
	}

	return "nil"
}
