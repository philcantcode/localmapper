package utils

import "time"

func Now() string {
	dt := time.Now()

	return dt.Format("02-01-2006 15:04:05")
}
