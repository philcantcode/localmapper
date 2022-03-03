package blueprint

type Vlan struct {
	ID          int    `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	HighIP      string `json:"HighIP"`
	LowIP       string `json:"LowIP"`
}
