package utils

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func Parse2dCSVFile(path string, delim string) map[string]string {

	var csv = make(map[string]string)

	file, err := os.Open(path)

	if err != nil {
		log.Fatal("Couldn't ParseFile ", path, err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		if strings.HasPrefix(scanner.Text(), "//") || len(scanner.Text()) == 0 {
			continue
		}

		elements := strings.Split(scanner.Text(), delim)
		csv[elements[0]] = elements[1]
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	file.Close()

	return csv
}

func ParseFile(path string) []string {
	input, err := ioutil.ReadFile("./" + path)
	if err != nil {
		log.Fatalln(err)
	}

	return strings.Split(string(input), "\n")
}

func ClearFile(path string) {

	if err := os.Truncate(path, 0); err != nil {
		log.Fatal("Failed to ClearFile: ", err)
	}
}

func CreateAndWriteFile(path string, contents string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	f.WriteString(contents)

	f.Close()
}

func DirExists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	}

	return false
}
