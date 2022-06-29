package utils

import (
	"fmt"
	"math"
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

func FormatDuration(dur time.Duration) string {
	output := ""

	if dur.Hours() >= 1 {
		hours := int64(math.Floor(dur.Hours()))
		output += fmt.Sprintf("%d hours ", hours)
		dur = time.Duration(int64(dur) - int64(time.Hour)*hours)
	}

	if dur.Minutes() >= 1 {
		mins := int64(math.Floor(dur.Minutes()))
		output += fmt.Sprintf("%d mins ", mins)
		dur = time.Duration(int64(dur) - int64(time.Minute)*mins)
	}

	if dur.Seconds() >= 1 {
		seconds := int64(math.Floor(dur.Seconds()))
		output += fmt.Sprintf("%d secs ", seconds)
		dur = time.Duration(int64(dur) - int64(time.Second)*seconds)
	}

	return output
}
