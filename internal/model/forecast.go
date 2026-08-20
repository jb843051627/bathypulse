package model

import (
	"sort"
	"time"
)

type ForecastPoint struct {
	At         time.Time `json:"at"`
	Expected   float64   `json:"expected"`
	Lower      float64   `json:"lower"`
	Upper      float64   `json:"upper"`
	Confidence float64   `json:"confidence"`
}

type ActivityForecast struct {
	StationID string          `json:"station_id"`
	Horizon   time.Duration   `json:"horizon"`
	Points    []ForecastPoint `json:"points"`
}

func (f ActivityForecast) Sort() ActivityForecast {
	points := append([]ForecastPoint(nil), f.Points...)
	sort.SliceStable(points, func(i, j int) bool { return points[i].At.Before(points[j].At) })
	f.Points = points
	return f
}

func (f ActivityForecast) At(at time.Time) (ForecastPoint, bool) {
	for _, point := range f.Points {
		if point.At.Equal(at) {
			return point, true
		}
	}
	return ForecastPoint{}, false
}

func (f ActivityForecast) Exceeds(point ForecastPoint, value float64) bool {
	return value > point.Upper && point.Confidence >= 0.5
}
