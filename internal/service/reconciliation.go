package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type Reconciliation struct {
	StationID string `json:"station_id"`
	Stored    int    `json:"stored"`
	Accepted  int    `json:"accepted"`
	Gap       int    `json:"gap"`
}

func (s *ObservatoryService) ReconcileSamples(ctx context.Context, stationID string, expected int) (Reconciliation, error) {
	if expected < 0 {
		return Reconciliation{}, model.ErrInvalid
	}
	_, _, err := s.samples.Aggregate(ctx, stationID, time.Time{}, s.clock.Now().Add(time.Minute))
	if err != nil {
		return Reconciliation{}, fmt.Errorf("reconcile samples: %w", err)
	}
	items, err := s.samples.Latest(ctx, stationID, expected+1)
	if err != nil {
		return Reconciliation{}, err
	}
	stored := len(items)
	return Reconciliation{StationID: stationID, Stored: stored, Accepted: expected, Gap: expected - stored}, nil
}

func (s *ObservatoryService) RequireComplete(reconciliation Reconciliation) error {
	if reconciliation.Gap != 0 {
		return fmt.Errorf("sample gap %d: %w", reconciliation.Gap, model.ErrConflict)
	}
	return nil
}
