package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bathypulse/internal/ingest"
	"github.com/jb843051627/bathypulse/internal/model"
)

func (s *ObservatoryService) SaveHealth(ctx context.Context, message ingest.HealthMessage) (model.StationHealth, error) {
	health, err := ingest.ValidateHealth(ctx, message)
	if err != nil {
		return model.StationHealth{}, err
	}
	if _, err := s.stations.Get(ctx, message.StationID); err != nil {
		return model.StationHealth{}, fmt.Errorf("health station: %w", err)
	}
	if err := s.db.Snapshots().Save(ctx, message.StationID, message.At, health); err != nil {
		return model.StationHealth{}, err
	}
	status := ingest.HealthState(health)
	station, err := s.stations.Get(ctx, message.StationID)
	if err != nil {
		return model.StationHealth{}, err
	}
	if status != station.Status && model.AllowedStationTransition(station.Status, status) {
		if err := s.stations.SetStatus(ctx, message.StationID, status, station.Version); err != nil {
			return model.StationHealth{}, err
		}
	}
	return health, nil
}

func (s *ObservatoryService) LatestHealth(ctx context.Context, stationID string) (*model.StationHealth, error) {
	return s.db.Snapshots().Latest(ctx, stationID)
}
