package report

import (
	"sort"

	"github.com/jb843051627/bathypulse/internal/model"
)

type NetworkReport struct {
	Stations []model.Station      `json:"stations"`
	Events   []model.SeismicEvent `json:"events"`
	Alerts   []model.Alert        `json:"alerts"`
	Health   map[string]float64   `json:"health"`
}

func BuildNetworkReport(stations []model.Station, events []model.SeismicEvent, alerts []model.Alert, health map[string]float64) NetworkReport {
	orderedStations := append([]model.Station(nil), stations...)
	orderedEvents := append([]model.SeismicEvent(nil), events...)
	orderedAlerts := append([]model.Alert(nil), alerts...)
	sort.SliceStable(orderedStations, func(i, j int) bool { return orderedStations[i].Code < orderedStations[j].Code })
	sort.SliceStable(orderedEvents, func(i, j int) bool { return orderedEvents[i].StartedAt.Before(orderedEvents[j].StartedAt) })
	sort.SliceStable(orderedAlerts, func(i, j int) bool { return orderedAlerts[i].CreatedAt.Before(orderedAlerts[j].CreatedAt) })
	copyHealth := make(map[string]float64, len(health))
	for id, value := range health {
		copyHealth[id] = value
	}
	return NetworkReport{Stations: orderedStations, Events: orderedEvents, Alerts: orderedAlerts, Health: copyHealth}
}

func OpenAlertCount(alerts []model.Alert) int {
	count := 0
	for _, alert := range alerts {
		if alert.IsActive() {
			count++
		}
	}
	return count
}
