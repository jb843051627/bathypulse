package model

import "time"

type StationHealth struct {
	StationID       string    `json:"station_id"`
	ObservedAt      time.Time `json:"observed_at"`
	BatteryPercent  float64   `json:"battery_percent"`
	SignalQuality   int       `json:"signal_quality"`
	ClockSkewMillis int64     `json:"clock_skew_millis"`
	SampleCount     int       `json:"sample_count"`
}

func (h StationHealth) Score() float64 {
	battery := h.BatteryPercent / 100
	signal := float64(h.SignalQuality) / 100
	skew := 1 - float64(abs64(h.ClockSkewMillis))/5000
	if skew < 0 {
		skew = 0
	}
	return (battery*0.4 + signal*0.4 + skew*0.2) * sampleFactor(h.SampleCount)
}

func (h StationHealth) Ready() bool {
	return h.BatteryPercent >= 15 && h.SignalQuality >= 50 && abs64(h.ClockSkewMillis) < 2000
}

func abs64(value int64) int64 {
	if value < 0 {
		return -value
	}
	return value
}

func sampleFactor(count int) float64 {
	if count <= 0 {
		return 0
	}
	if count >= 100 {
		return 1
	}
	return float64(count) / 100
}
