package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/report"
)

func (s *ObservatoryService) ReviewWaveform(ctx context.Context, waveformID string, review model.QualityReview) error {
	if review.WaveformID != waveformID || review.Score < 0 || review.Score > 1 {
		return model.ErrInvalid
	}
	waveform, err := s.waveforms.Get(ctx, waveformID)
	if err != nil {
		return fmt.Errorf("quality waveform: %w", err)
	}
	state := model.WaveformProcessed
	if !review.Accepted() {
		state = model.WaveformRejected
	}
	if err := s.waveforms.MarkProcessed(ctx, waveform.ID, state); err != nil {
		return err
	}
	s.metrics.Set("quality."+waveform.StationID, review.Score)
	return nil
}

func (s *ObservatoryService) QualitySummary(ctx context.Context, stationID string, reviews []model.QualityReview) (report.QualityReport, error) {
	if err := ctx.Err(); err != nil {
		return report.QualityReport{}, err
	}
	if _, err := s.stations.Get(ctx, stationID); err != nil {
		return report.QualityReport{}, err
	}
	return report.BuildQualityReport(stationID, reviews), nil
}
