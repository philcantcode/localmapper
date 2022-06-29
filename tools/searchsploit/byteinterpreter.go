package searchsploit

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/philcantcode/localmapper/system"
)

/*
	The results returned can be consecutive json objects {...}{....}{.....}
	which is technically malformed so doesn't parse correctly. We have to
	split the large string on the SEARCH key then rejoin removing whitespace
*/
var regex = regexp.MustCompile(`({[\s]+"SEARCH":)+`)
var header = "{ \"SEARCH\": "
var whitespace = regexp.MustCompile(`\s+`)

func ExtractExploitDB(resultBytes []byte) []ExploitDB {

	exploits := []ExploitDB{}
	regRes := regex.Split(string(resultBytes), -1)

	for _, searchTerm := range regRes {
		// Sometimes the searchterm is empty??
		if len(searchTerm) == 0 {
			continue
		}

		fixedSearchTerm := header + whitespace.ReplaceAllString(searchTerm, " ")

		exploitRes := ExploitDB{}
		err := json.Unmarshal([]byte(fixedSearchTerm), &exploitRes)
		system.Error(fmt.Sprintf("Couldn't unmarshal ExploitDB: %s", string(searchTerm)), err)

		if !exploitRes.isValid() {
			continue
		}

		exploits = append(exploits, exploitRes)
	}

	return exploits
}

/*
	Sometimes the results returned are empty, check
	to make sure they are not
*/
func (searchsploit ExploitDB) isValid() bool {

	if searchsploit.Search == "" {
		return false
	}

	if searchsploit.DbPathExploit == "" && searchsploit.DbPathShellcode == "" {
		return false
	}

	if len(searchsploit.ResultsExploit) == 0 && len(searchsploit.ResultsShellcode) == 0 {
		return false
	}

	return true
}
