package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// ErrorLog requires an err, non fatal
func ErrorLog(context string, e error, print bool) {
	if e != nil {
		if print {
			fmt.Println(context)
		}

		jsonErr, _ := json.Marshal(JsonLog{Type: "Error", DateTime: Now(), Context: context})
		AppendJson(string(jsonErr), Configs["JSON_LOG"])

		AppendLine("["+Now()+"] "+context, Configs["ERROR_LOG"])
	}
}

// ErrorFatal throws the error then exits
func ErrorFatal(context string, e error) {
	if e != nil {
		fmt.Println(context)

		AppendLine("["+Now()+"] "+context, Configs["ERROR_LOG"])

		jsonErr, _ := json.Marshal(JsonLog{Type: "Fatal", DateTime: Now(), Context: context})
		AppendJson(string(jsonErr), Configs["JSON_LOG"])

		log.Fatalf("%s: %v\n", context, e)
		os.Exit(0)
	}
}

// ErrorForceFatal throws the error then exits
func ErrorForceFatal(context string) {
	fmt.Println(context)

	jsonErr, _ := json.Marshal(JsonLog{Type: "Fatal", DateTime: Now(), Context: context})
	AppendJson(string(jsonErr), Configs["JSON_LOG"])

	AppendLine("["+Now()+"] "+context, Configs["ERROR_LOG"])
	os.Exit(0)
}

// ErrorContextLog logs a custom error, doesn't require a err to be passed
func ErrorContextLog(context string, print bool) {
	if print {
		fmt.Println(context)
	}

	jsonErr, _ := json.Marshal(JsonLog{Type: "Error", DateTime: Now(), Context: context})
	AppendJson(string(jsonErr), Configs["JSON_LOG"])

	AppendLine("["+Now()+"] "+context, Configs["ERROR_LOG"])
}

func ErrorDelim(context string) {
	AppendLine(context, Configs["ERROR_LOG"])
}
