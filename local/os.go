package local

import (
	"runtime"
)

// OSInfo returns
func OSInfo() map[string]string {
	var osInfo = make(map[string]string)

	osInfo["OS"] = osVersion()
	osInfo["Architecture"] = osArchitecture()

	return osInfo
}

func osVersion() string {
	os := runtime.GOOS
	switch os {
	case "windows":
		return "windows"
	case "darwin":
		return "mac"
	case "linux":
		return "linux"
	default:
		return os
	}
}

func osArchitecture() string {
	return runtime.GOARCH
}
