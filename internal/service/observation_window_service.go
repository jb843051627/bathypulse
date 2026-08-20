package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

func (s *ObservatoryService) ValidateWindow(ctx context.Context, window model.ObservationWindow) error {
	if !window.Valid() {
		return model.ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if window.Duration() > 6*time.Hour {
		return fmt.Errorf("observation window: %w", model.ErrInvalid)
	}
	if _, err := s.GetStation(ctx, window.StationID); err != nil {
		return err
	}
	return nil
}

func (s *ObservatoryService) CompleteWindow(ctx context.Context, stationID string, at time.Time) (model.ObservationWindow, error) {
	window, err := s.AnalysisWindow(ctx, stationID, at)
	if err != nil {
		return model.ObservationWindow{}, err
	}
	if err := s.ValidateWindow(ctx, window); err != nil {
		return model.ObservationWindow{}, err
	}
	return window, nil
}

func (s *ObservatoryService) WindowReady(ctx context.Context, stationID string, at time.Time) (bool, error) {
	window, err := s.CompleteWindow(ctx, stationID, at)
	if err != nil {
		return false, err
	}
	return window.Complete() && window.Activity() >= 1, nil
}

func (s *ObservatoryService) WindowState(ctx context.Context, stationID string, at time.Time) (string, error) {
	ready, err := s.WindowReady(ctx, stationID, at)
	if err != nil {
		return "", err
	}
	if ready {
		return "ready", nil
	}
	return "awaiting-data", nil
}
