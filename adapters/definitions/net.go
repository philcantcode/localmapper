package definitions

type Vlan struct {
	ID          int    `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	HighIP      string `json:"HighIP"`
	LowIP       string `json:"LowIP"`
	Tags        string `json:"Tags"`
}
