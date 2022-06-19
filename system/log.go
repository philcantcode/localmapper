package system

import (
	"fmt"
	"os"

	"github.com/philcantcode/localmapper/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const INT_TO_STRING_CONVERSION = "Could not convert from integer to string"

type LogEntry struct {
	ID       primitive.ObjectID `bson:"_id"`
	Type     string
	DateTime string
	Context  string
	Error    string
}

func Log(context string, debug bool) {
	log := LogEntry{Type: "Info", DateTime: utils.Now(), Context: context, Error: ""}
	INSERT_LogEntry(log)

	if debug {
		fmt.Printf("[%s] %s\n", log.Type, log.Context)
	}
}

func Error(context string, err error) {
	log := LogEntry{Type: "Error", DateTime: utils.Now(), Context: context, Error: "err.Error()"}

	if err != nil {
		INSERT_LogEntry(log)
		fmt.Printf("[%s] %s\n", log.Type, log.Context)
	}
}

func Fatal(context string, err error) {
	log := LogEntry{Type: "Fatal", DateTime: utils.Now(), Context: context, Error: "err.Error()"}

	if err != nil {
		fmt.Printf("[%s] %s\n", log.Type, log.Context)
		fmt.Println(err)
		INSERT_LogEntry(log)
		os.Exit(0)
	}
}

func Force(context string, fatal bool) {
	log := LogEntry{Type: "Error", DateTime: utils.Now(), Context: context, Error: ""}

	if fatal {
		log.Type = "Fatal"
		INSERT_LogEntry(log)
		fmt.Printf("[%s] %s\n", log.Type, log.Context)
		os.Exit(0)
	}

	INSERT_LogEntry(log)
	fmt.Printf("[%s] %s\n", log.Type, log.Context)
}
