package report

import (
	"sort"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type WindowReport struct {
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Waveforms  int       `json:"waveforms"`
	Events     int       `json:"events"`
	PeakEnergy float64   `json:"peak_energy"`
	StationID  string    `json:"station_id"`
}

func BuildWindowReport(waveforms []model.Waveform, events []model.SeismicEvent, start, end time.Time) WindowReport {
	ordered := append([]model.Waveform(nil), waveforms...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].CapturedAt.Before(ordered[j].CapturedAt) })
	energy := 0.0
	for _, waveform := range ordered {
		if !waveform.CapturedAt.Before(start) && waveform.CapturedAt.Before(end) {
			energy += waveform.Energy()
		}
	}
	stationID := ""
	if len(ordered) > 0 {
		stationID = ordered[0].StationID
	}
	return WindowReport{Start: start, End: end, Waveforms: len(ordered), Events: len(events), PeakEnergy: energy, StationID: stationID}
}
