package feeds

func GetAllWordlists() []Wordlist {
	return select_all_wordlists()
}
