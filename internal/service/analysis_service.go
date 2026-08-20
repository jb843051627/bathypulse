package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/report"
	"github.com/jb843051627/bathypulse/internal/validation"
)

func (s *ObservatoryService) StartAnalysis(ctx context.Context, run model.AnalysisRun) (model.AnalysisRun, error) {
	if err := validation.AnalysisRun(run); err != nil {
		return run, err
	}
	if _, err := s.stations.Get(ctx, run.StationID); err != nil {
		return run, fmt.Errorf("analysis station: %w", err)
	}
	if run.State == "" {
		run.State = "running"
	}
	if run.StartedAt.IsZero() {
		run.StartedAt = s.clock.Now()
	}
	if err := s.db.Analyses().Create(ctx, run); err != nil {
		return run, fmt.Errorf("create analysis: %w", err)
	}
	return run, nil
}

func (s *ObservatoryService) FinishAnalysis(ctx context.Context, id string, score float64) (model.AnalysisRun, error) {
	if score < 0 {
		return model.AnalysisRun{}, model.ErrInvalid
	}
	if err := s.db.Analyses().Complete(ctx, id, s.clock.Now(), score); err != nil {
		return model.AnalysisRun{}, fmt.Errorf("complete analysis: %w", err)
	}
	return s.db.Analyses().Get(ctx, id)
}

func (s *ObservatoryService) ListAnalyses(ctx context.Context, stationID string, limit int) ([]model.AnalysisRun, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.db.Analyses().List(ctx, stationID, limit)
}

func (s *ObservatoryService) AnalyzeWindow(ctx context.Context, stationID string, start, end time.Time) (report.AnalysisReport, error) {
	if err := validation.AnalysisWindow(start, end); err != nil {
		return report.AnalysisReport{}, err
	}
	waveforms, err := s.waveforms.ListByStation(ctx, stationID, 500)
	if err != nil {
		return report.AnalysisReport{}, err
	}
	events, err := s.events.FindOverlaps(ctx, stationID, start, end)
	if err != nil {
		return report.AnalysisReport{}, err
	}
	buckets := make([]model.AnalysisBucket, 0)
	for _, waveform := range waveforms {
		if waveform.CapturedAt.Before(start) || !waveform.CapturedAt.Before(end) {
			continue
		}
		buckets = append(buckets, model.AnalysisBucket{Start: waveform.CapturedAt, End: waveform.EndAt(), Mean: waveform.RMS, Peak: waveform.Peak, Count: 1})
	}
	run := model.AnalysisRun{ID: "adhoc-" + stationID, StationID: stationID, StartedAt: start, FinishedAt: s.clock.Now(), State: "complete", Samples: len(waveforms), Events: len(events), Score: report.AnalysisScore(buckets)}
	return report.BuildAnalysisReport(run, buckets, events), nil
}

func (s *ObservatoryService) AnalysisHealthy(ctx context.Context, stationID string) (bool, error) {
	items, err := s.ListAnalyses(ctx, stationID, 1)
	if err != nil {
		return false, err
	}
	if len(items) == 0 {
		return false, nil
	}
	return items[0].State == "complete" && items[0].Score >= 0.5, nil
}
