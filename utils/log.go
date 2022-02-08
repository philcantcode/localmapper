package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func Print(msg string) {
	fmt.Println(msg)
}

// Log resets after every new run
func Log(context string, print bool) {
	if print {
		fmt.Println(context)
	}

	AppendLine("["+Now()+"] "+context, Configs["RUNTIME_LOG"])
}

// PrintDump is intended to store the outputs of scan results
func PrintLog(context string) {
	AppendLine("["+Now()+"] "+context, Configs["PRINT_LOG"])
}

func FatalAlert(message string) {
	fmt.Println(message)
	os.Exit(0)
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func PrettyPrintToStr(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		return string(b)
	}

	return "ERROR PrettyPrintToStr"
}
