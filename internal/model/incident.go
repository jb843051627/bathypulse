package model

import "time"

type IncidentState string

const (
	IncidentOpen     IncidentState = "open"
	IncidentInvestig IncidentState = "investigating"
	IncidentClosed   IncidentState = "closed"
)

type Incident struct {
	ID        string        `json:"id"`
	StationID string        `json:"station_id"`
	AlertID   string        `json:"alert_id"`
	State     IncidentState `json:"state"`
	Summary   string        `json:"summary"`
	OpenedAt  time.Time     `json:"opened_at"`
	ClosedAt  *time.Time    `json:"closed_at,omitempty"`
	Revision  int           `json:"revision"`
}

func (i Incident) CanMoveTo(next IncidentState) bool {
	if i.State == IncidentOpen {
		return next == IncidentInvestig || next == IncidentClosed
	}
	return i.State == IncidentInvestig && next == IncidentClosed
}

func (i Incident) IsOpen() bool {
	return i.State != IncidentClosed
}
