package report

import (
	"sort"

	"github.com/jb843051627/bathypulse/internal/model"
)

type AnalysisReport struct {
	Run       model.AnalysisRun      `json:"run"`
	Buckets   []model.AnalysisBucket `json:"buckets"`
	TopEvents []model.SeismicEvent   `json:"top_events"`
}

func BuildAnalysisReport(run model.AnalysisRun, buckets []model.AnalysisBucket, events []model.SeismicEvent) AnalysisReport {
	orderedBuckets := model.SortBuckets(buckets)
	orderedEvents := append([]model.SeismicEvent(nil), events...)
	sort.SliceStable(orderedEvents, func(i, j int) bool { return orderedEvents[i].Magnitude > orderedEvents[j].Magnitude })
	if len(orderedEvents) > 10 {
		orderedEvents = orderedEvents[:10]
	}
	return AnalysisReport{Run: run, Buckets: orderedBuckets, TopEvents: orderedEvents}
}

func AnalysisScore(buckets []model.AnalysisBucket) float64 {
	if len(buckets) == 0 {
		return 0
	}
	total := 0.0
	for _, bucket := range buckets {
		total += model.BucketScore(bucket)
	}
	return total / float64(len(buckets))
}
