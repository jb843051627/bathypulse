package report

import (
	"sort"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type HealthReport struct {
	GeneratedAt time.Time             `json:"generated_at"`
	Stations    []model.StationHealth `json:"stations"`
	Average     float64               `json:"average"`
	Ready       int                   `json:"ready"`
}

func BuildHealthReport(at time.Time, items []model.StationHealth) HealthReport {
	ordered := append([]model.StationHealth(nil), items...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].Score() > ordered[j].Score() })
	result := HealthReport{GeneratedAt: at, Stations: ordered}
	for _, item := range ordered {
		result.Average += item.Score()
		if item.Ready() {
			result.Ready++
		}
	}
	if len(ordered) > 0 {
		result.Average /= float64(len(ordered))
	}
	return result
}

func HealthGrade(score float64) string {
	if score >= 0.85 {
		return "green"
	}
	if score >= 0.6 {
		return "amber"
	}
	return "red"
}
