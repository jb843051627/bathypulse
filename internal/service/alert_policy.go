package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bathypulse/internal/model"
)

func (s *ObservatoryService) EvaluateWaveform(ctx context.Context, waveform model.Waveform) (*model.Alert, error) {
	threshold, err := s.EventRatioThreshold(ctx)
	if err != nil {
		return nil, err
	}
	if waveform.RMS <= 0 || waveform.Peak/waveform.RMS < threshold {
		return nil, nil
	}
	event := model.SeismicEvent{ID: "policy-" + waveform.ID, StationID: waveform.StationID, Kind: "threshold", StartedAt: waveform.CapturedAt, EndedAt: waveform.EndAt(), Magnitude: waveform.Peak, Confidence: 0.7, State: model.EventOpen, Revision: 1, WaveformCount: 1}
	if err := s.events.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("policy event: %w", err)
	}
	return s.RaiseAlert(ctx, event, model.AlertWatch, "waveform ratio exceeded policy")
}

func (s *ObservatoryService) EscalatePending(ctx context.Context, stationID string) (int, error) {
	alerts, err := s.alerts.List(ctx, stationID, model.AlertPending)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, alert := range alerts {
		if alert.NeedsEscalation() {
			if err := s.alerts.SetLevel(ctx, alert.ID, model.AlertCritical, alert.Revision); err != nil {
				return count, err
			}
			count++
		}
	}
	return count, nil
}
