package unifi

import (
	"time"
)

type UnClients struct {
	Count float64 `json:"count"`
	Data  []struct {
		ConnectedAt    time.Time `json:"connectedAt"`
		ID             string    `json:"id"`
		IpAddress      string    `json:"ipAddress"`
		MacAddress     string    `json:"macAddress"`
		Name           string    `json:"name"`
		Type           string    `json:"type"`
		UplinkDeviceID string    `json:"uplinkDeviceId"`
	} `json:"data"`
	Limit      float64 `json:"limit"`
	Offset     float64 `json:"offset"`
	TotalCount float64 `json:"totalCount"`
}
