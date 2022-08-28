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

func insert_wordlist(wordlist Wordlist) {
	stmt, err := system.Con.Prepare("INSERT INTO `Wordlists` (`label`, `description`, `path`, `type`, `size`) VALUES (?, ?, ?, ?, ?);")
	system.Fatal("Couldn't INSERT Into Wordlists", err)

	_, err = stmt.Exec(wordlist.Label, wordlist.Description, wordlist.Path, wordlist.Type, wordlist.Size)
	system.Fatal("Results error from INSERT Wordlists", err)

	stmt.Close()
}
