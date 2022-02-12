package utils

import (
	"fmt"
	"log"
	"os"
)

func ErrorHandle(context string, e error, print bool) {
	if e != nil {
		if print {
			fmt.Println(context)
		}

		AppendLine("["+Now()+"] "+context, Configs["ERROR_LOG"])

		log.Fatalf("%s: %v\n", context, e)
	}
}

func FatalErrorHandle(context string, e error) {
	if e != nil {
		fmt.Println(context)

		AppendLine("["+Now()+"] "+context, Configs["ERROR_LOG"])

		log.Fatalf("%s: %v\n", context, e)
		os.Exit(0)
	}
}

func ErrorHandleLog(context string, print bool) {
	if print {
		fmt.Println(context)
	}

	AppendLine("["+Now()+"] "+context, Configs["ERROR_LOG"])
}

func ErrorDelim(context string) {
	AppendLine(context, Configs["ERROR_LOG"])
}
