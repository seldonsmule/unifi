package unifi

type UnDevices struct {
	Count float64 `json:"count"`
	Data  []struct {
		Features   []string `json:"features"`
		ID         string   `json:"id"`
		Interfaces []string `json:"interfaces"`
		IpAddress  string   `json:"ipAddress"`
		MacAddress string   `json:"macAddress"`
		Model      string   `json:"model"`
		Name       string   `json:"name"`
		State      string   `json:"state"`
	} `json:"data"`
	Limit      float64 `json:"limit"`
	Offset     float64 `json:"offset"`
	TotalCount float64 `json:"totalCount"`
}
