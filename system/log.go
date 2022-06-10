package system

import (
	"fmt"
	"os"

	"github.com/philcantcode/localmapper/utils"
)

const INT_TO_STRING_CONVERSION = "Could not convert from integer to string"

type LogEntry struct {
	Type     string
	DateTime string
	Context  string
	Error    string
}

func Log(context string, debug bool) {
	log := LogEntry{Type: "Info", DateTime: utils.Now(), Context: context, Error: ""}

	if debug {
		fmt.Printf("> %+v\n", log)
	}
}

func Error(context string, err error) {
	log := LogEntry{Type: "Error", DateTime: utils.Now(), Context: context, Error: "err.Error()"}

	if err != nil {
		fmt.Printf("> %+v\n", log)
	}
}

func Fatal(context string, err error) {
	log := LogEntry{Type: "Fatal", DateTime: utils.Now(), Context: context, Error: "err.Error()"}

	if err != nil {
		fmt.Printf("> %+v\n", log)
		os.Exit(0)
	}
}

func Force(context string, fatal bool) {
	log := LogEntry{Type: "Error", DateTime: utils.Now(), Context: context, Error: ""}

	if fatal {
		log.Type = "Fatal"
		fmt.Printf("> %+v\n", log)
		os.Exit(0)
	}

	fmt.Printf("> %+v\n", log)
}
