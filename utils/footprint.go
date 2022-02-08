package utils

import "runtime"

const (
	WINDOWS = "windows"
	LINUX   = "linux"
	OSX     = "darwin"
)

// IdentifyOS returns the OS version (windows, linux, darwin)
func IdentifyOS() string {
	return runtime.GOOS
}
