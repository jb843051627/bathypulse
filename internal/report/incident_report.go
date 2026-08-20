package report

import (
	"sort"

	"github.com/jb843051627/bathypulse/internal/model"
)

type IncidentReport struct {
	StationID string           `json:"station_id"`
	Open      []model.Incident `json:"open"`
	Alerts    []model.Alert    `json:"alerts"`
	Score     float64          `json:"score"`
}

func BuildIncidentReport(stationID string, incidents []model.Incident, alerts []model.Alert) IncidentReport {
	open := append([]model.Incident(nil), incidents...)
	orderedAlerts := append([]model.Alert(nil), alerts...)
	sort.SliceStable(open, func(i, j int) bool { return open[i].OpenedAt.Before(open[j].OpenedAt) })
	sort.SliceStable(orderedAlerts, func(i, j int) bool { return orderedAlerts[i].CreatedAt.Before(orderedAlerts[j].CreatedAt) })
	score := float64(len(open))*2 + float64(OpenAlertCount(orderedAlerts))
	return IncidentReport{StationID: stationID, Open: open, Alerts: orderedAlerts, Score: score}
}

func IncidentSeverity(item model.Incident, alerts []model.Alert) model.AlertLevel {
	for _, alert := range alerts {
		if alert.ID == item.AlertID {
			return alert.Level
		}
	}
	return model.AlertInfo
}
