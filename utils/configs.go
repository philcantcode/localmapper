package utils

import (
	"fmt"
	"os"
	"strings"
)

var Configs = make(map[string]string)

func LoadGlobalConfigs() {
	Configs = Parse2dCSVFile("res/conf.txt", ", ")

	// Create log files if they don't exist
	CreateFile(Configs["RUNTIME_LOG"])
	CreateFile(Configs["ERROR_LOG"])
	CreateFile(Configs["PRINT_LOG"])

	// Clear the logfiles that refersh on every run
	ClearFile(Configs["RUNTIME_LOG"])
	ClearFile(Configs["PRINT_LOG"])

	// Add space to the log files
	ErrorDelim("\n##########################| " + Now() + " |##########################")
}

func SaveGlobalConfigs() {
	lines := ParseFile("res/conf.txt")

	fmt.Printf("%v\n", lines)
	os.Exit(0)

	ClearFile("res/conf.txt")

	for i, line := range lines {
		for k := range Configs {
			if strings.HasPrefix(line, k) {
				lines[i] = k + ", " + Configs[k]
				AppendLine("res/conf.txt", lines[i])
			}
		}
	}
}
