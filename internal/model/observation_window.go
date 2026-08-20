package model

import "time"

type ObservationWindow struct {
	ID        string    `json:"id"`
	StationID string    `json:"station_id"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	State     string    `json:"state"`
	Samples   int       `json:"samples"`
	Events    int       `json:"events"`
}

func (w ObservationWindow) Valid() bool {
	return w.ID != "" && w.StationID != "" && !w.Start.IsZero() && w.End.After(w.Start)
}

func (w ObservationWindow) Duration() time.Duration {
	if !w.End.After(w.Start) {
		return 0
	}
	return w.End.Sub(w.Start)
}

func (w ObservationWindow) Complete() bool {
	return w.State == "complete" && w.Samples > 0
}

func (w ObservationWindow) Activity() float64 {
	if w.Duration() <= 0 {
		return 0
	}
	return float64(w.Samples+w.Events*10) / w.Duration().Minutes()
}
