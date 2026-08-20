package report

import (
	"math"

	"github.com/jb843051627/bathypulse/internal/model"
)

type Anomaly struct {
	StationID string  `json:"station_id"`
	Mean      float64 `json:"mean"`
	Spread    float64 `json:"spread"`
	Score     float64 `json:"score"`
}

func Analyze(samples []model.Sample) Anomaly {
	if len(samples) == 0 {
		return Anomaly{}
	}
	var sum float64
	for _, sample := range samples {
		sum += sample.Value
	}
	mean := sum / float64(len(samples))
	var variance float64
	for _, sample := range samples {
		delta := sample.Value - mean
		variance += delta * delta
	}
	spread := math.Sqrt(variance / float64(len(samples)))
	score := 0.0
	if spread > 0 {
		score = math.Abs(mean) / spread
	}
	return Anomaly{StationID: samples[0].StationID, Mean: mean, Spread: spread, Score: score}
}
