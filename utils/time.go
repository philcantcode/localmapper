package utils

import (
	"time"
)

func Now() string {
	dt := time.Now()

	return dt.Format("02-01-2006 15:04:05")
}

type DateTime struct {
	DDMMYYYY string
	HHMMSS   string
	DateTime string
}

const DTF_DDMMYYYY = "02-01-2006"
const DTF_HHMMSS = "15:04:05"
const DTF_DateTime = "02-01-2006 15:04:05"

func GetDateTime() DateTime {
	dt := time.Now()
	dts := DateTime{}

	dts.DDMMYYYY = dt.Format(DTF_DDMMYYYY)
	dts.HHMMSS = dt.Format(DTF_HHMMSS)
	dts.DateTime = dt.Format(DTF_DateTime)

	return dts
}
