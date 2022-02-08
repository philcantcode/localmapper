package utils

import "time"

func Now() string {
	dt := time.Now()

	return dt.Format("01-02-2006 15:04:05")
}
