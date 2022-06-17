package utils

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(str interface{}) {
	b, err := json.MarshalIndent(str, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(b))
}
