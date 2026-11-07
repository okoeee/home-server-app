package model

import "time"

type ServerMetrics struct {
	FetchedAt             time.Time
	CPUUsagePercent       float64
	MemoryUsagePercent    float64
	CPUTemperatureCelsius *float64
}
