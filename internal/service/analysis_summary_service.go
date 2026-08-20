package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type AnalysisSummary struct {
	StationID      string
	Runs           int
	Complete       int
	AverageScore   float64
	LatestFinished time.Time
}

func (s *ObservatoryService) SummarizeAnalyses(ctx context.Context, stationID string) (AnalysisSummary, error) {
	if stationID == "" {
		return AnalysisSummary{}, model.ErrInvalid
	}
	items, err := s.ListAnalyses(ctx, stationID, 100)
	if err != nil {
		return AnalysisSummary{}, fmt.Errorf("analysis summary: %w", err)
	}
	result := AnalysisSummary{StationID: stationID, Runs: len(items)}
	for _, item := range items {
		if item.State == "complete" {
			result.Complete++
		}
		result.AverageScore += item.Score
		if item.FinishedAt.After(result.LatestFinished) {
			result.LatestFinished = item.FinishedAt
		}
	}
	if result.Runs > 0 {
		result.AverageScore /= float64(result.Runs)
	}
	return result, nil
}

func (s *ObservatoryService) AnalysisWindow(ctx context.Context, stationID string, at time.Time) (model.ObservationWindow, error) {
	start, end, err := s.WindowForStation(ctx, stationID, at)
	if err != nil {
		return model.ObservationWindow{}, err
	}
	report, err := s.BuildWindow(ctx, stationID, start, end)
	if err != nil {
		return model.ObservationWindow{}, err
	}
	return model.ObservationWindow{ID: fmt.Sprintf("window-%s-%d", stationID, start.Unix()), StationID: stationID, Start: start, End: end, State: "complete", Samples: report.Waveforms, Events: report.Events}, nil
}

func (s *ObservatoryService) WindowActivity(ctx context.Context, stationID string, at time.Time) (float64, error) {
	window, err := s.AnalysisWindow(ctx, stationID, at)
	if err != nil {
		return 0, err
	}
	return window.Activity(), nil
}
