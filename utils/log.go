package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
)

type JsonLog struct {
	Type     string
	DateTime string
	Context  string
}

func Print(context string) {
	fmt.Println(context)

	jsonErr, _ := json.Marshal(JsonLog{Type: "Info", DateTime: Now(), Context: context})
	AppendJson(string(jsonErr), Configs["JSON_LOG"])
}

// Log resets after every new run
func Log(context string, print bool) {
	if print {
		fmt.Println(context)
	}

	jsonErr, _ := json.Marshal(JsonLog{Type: "Info", DateTime: Now(), Context: context})

	AppendJson(string(jsonErr), Configs["JSON_LOG"])
	AppendLine("[i]["+Now()+"] "+context, Configs["RUNTIME_LOG"])
}

// PrintDump is intended to store the outputs of scan results
func PrintLog(context string) {

	jsonErr, _ := json.Marshal(JsonLog{Type: "Info", DateTime: Now(), Context: context})
	AppendJson(string(jsonErr), Configs["JSON_LOG"])

	AppendLine("["+Now()+"] "+context, Configs["PRINT_LOG"])
}

func FatalAlert(context string) {
	fmt.Println(context)

	jsonErr, _ := json.Marshal(JsonLog{Type: "Fatal", DateTime: Now(), Context: context})
	AppendJson(string(jsonErr), Configs["JSON_LOG"])

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

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
