package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/jb843051627/bathypulse/internal/model"
)

type QualityBatchResult struct {
	StationID string
	Accepted  int
	Rejected  int
	Reviewed  int
	Scores    []float64
}

func (s *ObservatoryService) ReviewQualityBatch(ctx context.Context, stationID string, reviews []model.QualityReview) (QualityBatchResult, error) {
	if stationID == "" || len(reviews) == 0 {
		return QualityBatchResult{}, model.ErrInvalid
	}
	if _, err := s.stations.Get(ctx, stationID); err != nil {
		return QualityBatchResult{}, err
	}
	result := QualityBatchResult{StationID: stationID, Scores: make([]float64, 0, len(reviews))}
	for _, review := range reviews {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		if review.WaveformID == "" || review.Score < 0 || review.Score > 1 {
			return result, fmt.Errorf("quality review: %w", model.ErrInvalid)
		}
		result.Scores = append(result.Scores, review.Score)
		if review.Accepted() {
			result.Accepted++
		} else {
			result.Rejected++
		}
		if review.NeedsOperator() {
			result.Reviewed++
		}
	}
	sort.Float64s(result.Scores)
	return result, nil
}

func (r QualityBatchResult) Average() float64 {
	if len(r.Scores) == 0 {
		return 0
	}
	total := 0.0
	for _, score := range r.Scores {
		total += score
	}
	return total / float64(len(r.Scores))
}

func (r QualityBatchResult) Complete() bool {
	return r.Accepted > 0 && r.Rejected == 0 && r.Reviewed == 0
}

func (s *ObservatoryService) QualityGate(ctx context.Context, stationID string, reviews []model.QualityReview) error {
	result, err := s.ReviewQualityBatch(ctx, stationID, reviews)
	if err != nil {
		return err
	}
	if !result.Complete() || result.Average() < 0.75 {
		return fmt.Errorf("quality gate: %w", model.ErrConflict)
	}
	return nil
}

func (s *ObservatoryService) QualityMetrics(ctx context.Context, stationID string, reviews []model.QualityReview) (map[string]float64, error) {
	result, err := s.ReviewQualityBatch(ctx, stationID, reviews)
	if err != nil {
		return nil, err
	}
	metrics := map[string]float64{"average": result.Average(), "accepted": float64(result.Accepted), "rejected": float64(result.Rejected), "reviewed": float64(result.Reviewed)}
	return metrics, nil
}
