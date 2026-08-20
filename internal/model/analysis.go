package model

import (
	"sort"
	"time"
)

type AnalysisRun struct {
	ID         string    `json:"id"`
	StationID  string    `json:"station_id"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	State      string    `json:"state"`
	Samples    int       `json:"samples"`
	Events     int       `json:"events"`
	Score      float64   `json:"score"`
}

type AnalysisBucket struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Mean  float64   `json:"mean"`
	Peak  float64   `json:"peak"`
	Count int       `json:"count"`
}

func (r AnalysisRun) Complete(at time.Time, score float64) AnalysisRun {
	r.FinishedAt = at
	r.State = "complete"
	r.Score = score
	return r
}

func (r AnalysisRun) Duration() time.Duration {
	if r.FinishedAt.Before(r.StartedAt) {
		return 0
	}
	return r.FinishedAt.Sub(r.StartedAt)
}

func SortBuckets(buckets []AnalysisBucket) []AnalysisBucket {
	result := append([]AnalysisBucket(nil), buckets...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].Start.Before(result[j].Start) })
	return result
}

func BucketScore(bucket AnalysisBucket) float64 {
	if bucket.Count == 0 {
		return 0
	}
	return bucket.Peak + bucket.Mean*0.25 + float64(bucket.Count)*0.01
}
