package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bathypulse/internal/model"
)

func (s *ObservatoryService) ReserveSequences(ctx context.Context, stationID string, count int64) (int64, error) {
	if stationID == "" || count < 1 || count > 100000 {
		return 0, model.ErrInvalid
	}
	if _, err := s.stations.Get(ctx, stationID); err != nil {
		return 0, fmt.Errorf("sequence station: %w", err)
	}
	return s.db.Sequences().Reserve(ctx, stationID, count)
}

func (s *ObservatoryService) CurrentSequence(ctx context.Context, stationID string) (int64, error) {
	if stationID == "" {
		return 0, model.ErrInvalid
	}
	return s.db.Sequences().Current(ctx, stationID)
}

func (s *ObservatoryService) SequenceHealthy(ctx context.Context, stationID string) (bool, error) {
	value, err := s.CurrentSequence(ctx, stationID)
	if err != nil {
		return false, err
	}
	return value > 0, nil
}
