package unifi

import (
	"time"
)

type UnDeviceStats struct {
	CpuUtilizationPct    float64   `json:"cpuUtilizationPct"`
	Interfaces           struct{}  `json:"interfaces"`
	LastHeartbeatAt      time.Time `json:"lastHeartbeatAt"`
	LoadAverage15Min     float64   `json:"loadAverage15Min"`
	LoadAverage1Min      float64   `json:"loadAverage1Min"`
	LoadAverage5Min      float64   `json:"loadAverage5Min"`
	MemoryUtilizationPct float64   `json:"memoryUtilizationPct"`
	NextHeartbeatAt      time.Time `json:"nextHeartbeatAt"`
	Uplink               struct {
		RxRateBps float64 `json:"rxRateBps"`
		TxRateBps float64 `json:"txRateBps"`
	} `json:"uplink"`
	UptimeSec float64 `json:"uptimeSec"`
}
