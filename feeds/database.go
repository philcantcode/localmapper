package feeds

import (
	"github.com/philcantcode/localmapper/system"
)

func select_all_wordlists() []Wordlist {
	stmt, err := system.Con.Query("SELECT `id`, `label`, `description`, `path`, `type`, `size` FROM `Wordlists` ORDER BY `type`;")
	system.Fatal("Couldn't SELECT_ALL from Wordlists", err)

	var allWordlists []Wordlist

	for stmt.Next() {
		wordlist := Wordlist{}

		stmt.Scan(&wordlist.ID, &wordlist.Label, &wordlist.Description, &wordlist.Path, &wordlist.Type, &wordlist.Size)
		allWordlists = append(allWordlists, wordlist)
	}

	stmt.Close()
	return allWordlists
}
