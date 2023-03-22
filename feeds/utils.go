package feeds

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ParseWordlistFolder() {
	searchDir := "/localmapper/wordlists/usernames"

	fileList := make([]string, 0)
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return err
	})

	data := []Wordlist{}

	for _, file := range fileList {
		wl := Wordlist{}
		fileStat, err := os.Stat(file)

		if fileStat.IsDir() {
			continue
		}

		if err != nil {
			fmt.Println(err)
		}

		wl.Label = strings.Title(strings.ReplaceAll(strings.Split(fileStat.Name(), ".")[0], "-", " ")) + " Username List"
		wl.Description = "List of usernames called " + strings.Title(strings.ReplaceAll(strings.Split(fileStat.Name(), ".")[0], "-", " "))
		wl.Path = strings.Split(file, "/localmapper/")[1]
		wl.Type = "Username"
		wl.Size = fileStat.Size()

		data = append(data, wl)

		fmt.Printf("%+v\n", wl)
	}

	a := ""
	fmt.Scan(&a)

	for _, wl := range data {
		insert_wordlist(wl)
	}

}
