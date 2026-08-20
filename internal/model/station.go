package model

import "time"

func (s Station) IsUsable(at time.Time) bool {
	if s.Status == StationRepair || s.Status == StationOffline {
		return false
	}
	if s.LastSeen.IsZero() {
		return false
	}
	return at.Sub(s.LastSeen) <= 15*time.Minute
}

func (s Station) Coordinates() [2]float64 {
	return [2]float64{s.Latitude, s.Longitude}
}

func (s Station) DisplayName() string {
	if s.Code == "" {
		return s.ID
	}
	return s.Code + " @ " + s.Basin
}
