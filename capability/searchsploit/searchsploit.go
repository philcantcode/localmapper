package searchsploit

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/philcantcode/localmapper/system"
	"go.mongodb.org/mongo-driver/bson"
)

/*
	The results returned can be consecutive json objects {...}{....}{.....}
	which is technically malformed so doesn't parse correctly. We have to
	split the large string on the SEARCH key then rejoin removing whitespace
*/
var regex = regexp.MustCompile(`({[\s]+"SEARCH":)+`)
var header = "{ \"SEARCH\": "
var whitespace = regexp.MustCompile(`\s+`)

func ProcessResults(res []byte) []ExploitDB {

	exploits := []ExploitDB{}
	regRes := regex.Split(string(res), -1)

	for _, searchTerm := range regRes {
		// Sometimes the searchterm is empty??
		if len(searchTerm) == 0 {
			continue
		}

		fixedSearchTerm := header + whitespace.ReplaceAllString(searchTerm, " ")

		exploitRes := ExploitDB{}
		err := json.Unmarshal([]byte(fixedSearchTerm), &exploitRes)
		system.Error(fmt.Sprintf("Couldn't unmarshal ExploitDB: %s", string(searchTerm)), err)

		exploits = append(exploits, exploitRes)
	}

	return exploits
}

func Test() {
	test := `{ "SEARCH": "dropbear sshd",      "DB_PATH_EXPLOIT": "/localmapper/exploitdb",    "RESULTS_EXPLOIT": [  {"Title":"DropBearSSHD 2015.71 - Command Injection","EDB-ID":"40119","Date":"1970-01-01","Author":"tintinweb","Type":"remote","Platform":"linux","Path":"/localmapper/exploitdb/exploits/linux/remote/40119.md"}        ],      "DB_PATH_SHELLCODE": "/localmapper/exploitdb",  "RESULTS_SHELLCODE": [ ]}{ "SEARCH": "dropbear sshd",      "DB_PATH_EXPLOIT": "/localmapper/exploitdb",    "RESULTS_EXPLOIT": [  {"Title":"DropBearSSHD 2015.71 - Command Injection","EDB-ID":"40119","Date":"1970-01-01","Author":"tintinweb","Type":"remote","Platform":"linux","Path":"/localmapper/exploitdb/exploits/linux/remote/40119.md"}        ],      "DB_PATH_SHELLCODE": "/localmapper/exploitdb",  "RESULTS_SHELLCODE": [ ]}{ "SEARCH": "dropbear sshd",      "DB_PATH_EXPLOIT": "/localmapper/exploitdb",    "RESULTS_EXPLOIT": [  {"Title":"DropBearSSHD 2015.71 - Command Injection","EDB-ID":"40119","Date":"1970-01-01","Author":"tintinweb","Type":"remote","Platform":"linux","Path":"/localmapper/exploitdb/exploits/linux/remote/40119.md"}        ],      "DB_PATH_SHELLCODE": "/localmapper/exploitdb",  "RESULTS_SHELLCODE": [ ]}`

	fmt.Println("Trying: " + test)

	regRes := regex.Split(string(test), -1)

	for i, searchTerm := range regRes {
		fixedSearchTerm := header + whitespace.ReplaceAllString(searchTerm, " ")

		fmt.Printf("[%d]: %s\n\n", i, fixedSearchTerm)

		exploitRes := ExploitDB{}
		err := json.Unmarshal([]byte(fixedSearchTerm), &exploitRes)

		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
	}

	os.Exit(0)
}

func (result ExploitDB) StoreResults() string {
	if result.Search == "" {
		return ""
	}

	if result.DbPathExploit == "" && result.DbPathShellcode == "" {
		return ""
	}

	if len(result.ResultsExploit) == 0 && len(result.ResultsShellcode) == 0 {
		return ""
	}

	// Don't input where the exact search has been done before
	if len(select_searchsploit(bson.M{"search": result.Search}, bson.M{})) != 0 {
		delete_searchsploit(bson.M{"search": result.Search})
	}

	return insert_searchsploit(result)
}
