package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

func (s *ObservatoryService) DefaultPolicy() model.ObservatoryPolicy {
	return model.ObservatoryPolicy{ID: "default", Name: "deep-water-observation", RatioThreshold: s.PolicySnapshot(), MaxWindow: 30 * time.Minute, ReviewAfter: 10 * time.Minute, Enabled: true}
}

func (s *ObservatoryService) CheckWaveformPolicy(ctx context.Context, waveform model.Waveform) error {
	policy, err := s.EventRatioThreshold(ctx)
	if err != nil {
		return err
	}
	current := s.DefaultPolicy()
	current.RatioThreshold = policy
	if !current.WindowAllowed(time.Duration(waveform.DurationMs) * time.Millisecond) {
		return fmt.Errorf("waveform window policy: %w", model.ErrInvalid)
	}
	if !current.StrongRatio(waveform.Peak, waveform.RMS) {
		return fmt.Errorf("waveform strength policy: %w", model.ErrInvalid)
	}
	return nil
}

func (s *ObservatoryService) PolicyNeedsReview(ctx context.Context, openedAt time.Time) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	return s.DefaultPolicy().ShouldReview(s.clock.Now().Sub(openedAt)), nil
}
