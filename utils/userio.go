package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func UserStringInput(prompt string) string {
	fmt.Println(prompt)
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		Log("Couldn't read user input from console", false)
	}

	return strings.TrimSpace(input)
}

func UserPrint(message string) {
	fmt.Println(message)
}
