package unifi

type UnSites struct {
	Count float64 `json:"count"`
	Data  []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
	Limit      float64 `json:"limit"`
	Offset     float64 `json:"offset"`
	TotalCount float64 `json:"totalCount"`
}
