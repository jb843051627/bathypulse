package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/codec"
	"github.com/jb843051627/bathypulse/internal/model"
)

func (s *ObservatoryService) SaveReplay(ctx context.Context, id string, record codec.ReplayRecord) error {
	if id == "" || record.StationID == "" || record.At.IsZero() {
		return fmt.Errorf("replay identity: %w", model.ErrInvalid)
	}
	if _, err := s.stations.Get(ctx, record.StationID); err != nil {
		return err
	}
	return s.db.Replays().Save(ctx, id, record)
}

func (s *ObservatoryService) LoadReplay(ctx context.Context, stationID string, start, end time.Time) ([]codec.ReplayRecord, error) {
	if start.IsZero() || !end.After(start) {
		return nil, model.ErrInvalid
	}
	return s.db.Replays().List(ctx, stationID, start, end)
}

func (s *ObservatoryService) PruneReplay(ctx context.Context, stationID string, before time.Time) (int64, error) {
	if before.IsZero() {
		return 0, model.ErrInvalid
	}
	return s.db.Replays().DeleteBefore(ctx, stationID, before)
}
