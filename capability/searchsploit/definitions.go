package searchsploit

type ExploitDB struct {
	Search         string `json:"SEARCH"`
	DbPathExploit  string `json:"DB_PATH_EXPLOIT"`
	ResultsExploit []struct {
		Title    string `json:"Title"`
		EDBID    string `json:"EDB-ID"`
		Date     string `json:"Date"`
		Author   string `json:"Author"`
		Type     string `json:"Type"`
		Platform string `json:"Platform"`
		Path     string `json:"Path"`
	} `json:"RESULTS_EXPLOIT"`
	DbPathShellcode  string `json:"DB_PATH_SHELLCODE"`
	ResultsShellcode []struct {
		Title    string `json:"Title"`
		EDBID    string `json:"EDB-ID"`
		Date     string `json:"Date"`
		Author   string `json:"Author"`
		Type     string `json:"Type"`
		Platform string `json:"Platform"`
		Path     string `json:"Path"`
	} `json:"RESULTS_SHELLCODE"`
}
