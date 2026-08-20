package report

import (
	"sort"

	"github.com/jb843051627/bathypulse/internal/model"
)

type CalibrationReport struct {
	StationID string                   `json:"station_id"`
	Version   int                      `json:"version"`
	Samples   int                      `json:"samples"`
	MeanError float64                  `json:"mean_error"`
	Points    []model.CalibrationPoint `json:"points"`
}

func BuildCalibrationReport(profile model.CalibrationProfile, values []float64) CalibrationReport {
	points := append([]model.CalibrationPoint(nil), profile.Points...)
	sort.Slice(points, func(i, j int) bool { return points[i].Input < points[j].Input })
	var errorSum float64
	for _, value := range values {
		errorSum += profile.Apply(value) - value
	}
	mean := 0.0
	if len(values) > 0 {
		mean = errorSum / float64(len(values))
	}
	return CalibrationReport{StationID: profile.StationID, Version: profile.Version, Samples: len(values), MeanError: mean, Points: points}
}
