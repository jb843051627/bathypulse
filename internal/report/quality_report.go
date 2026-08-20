package report

import (
	"sort"

	"github.com/jb843051627/bathypulse/internal/model"
)

type QualityReport struct {
	StationID string  `json:"station_id"`
	Total     int     `json:"total"`
	Accepted  int     `json:"accepted"`
	Reviewed  int     `json:"reviewed"`
	Discarded int     `json:"discarded"`
	MeanScore float64 `json:"mean_score"`
}

func BuildQualityReport(stationID string, reviews []model.QualityReview) QualityReport {
	ordered := append([]model.QualityReview(nil), reviews...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].WaveformID < ordered[j].WaveformID })
	result := QualityReport{StationID: stationID, Total: len(ordered)}
	for _, review := range ordered {
		if review.Accepted() {
			result.Accepted++
		}
		if review.NeedsOperator() {
			result.Reviewed++
		}
		if review.Flag == model.QualityDiscarded {
			result.Discarded++
		}
		result.MeanScore += review.Score
	}
	if result.Total > 0 {
		result.MeanScore /= float64(result.Total)
	}
	return result
}
