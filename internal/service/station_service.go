package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/validation"
)

func (s *ObservatoryService) RegisterStation(ctx context.Context, station model.Station) (*model.Station, error) {
	if err := validation.Station(station); err != nil {
		return nil, err
	}
	if station.Status == "" {
		station.Status = model.StationActive
	}
	if station.Version == 0 {
		station.Version = 1
	}
	if err := s.stations.Create(ctx, station); err != nil {
		return nil, fmt.Errorf("register station: %w", err)
	}
	return &station, nil
}

func (s *ObservatoryService) GetStation(ctx context.Context, id string) (*model.Station, error) {
	return s.stations.Get(ctx, id)
}

func (s *ObservatoryService) ListStations(ctx context.Context, status model.StationStatus) ([]model.Station, error) {
	return s.stations.List(ctx, status)
}

func (s *ObservatoryService) UpdateHeartbeat(ctx context.Context, id string, latency float64) error {
	station, err := s.stations.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("load station heartbeat: %w", err)
	}
	if err := s.stations.UpdateHeartbeat(ctx, id, s.clock.Now(), station.Version); err != nil {
		return fmt.Errorf("store station heartbeat: %w", err)
	}
	s.metrics.ObserveStation(id, latency)
	return nil
}

func (s *ObservatoryService) ChangeStationStatus(ctx context.Context, id string, next model.StationStatus) error {
	station, err := s.stations.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("load station status: %w", err)
	}
	if !model.AllowedStationTransition(station.Status, next) {
		return fmt.Errorf("station transition %s to %s: %w", station.Status, next, model.ErrConflict)
	}
	if err := s.stations.SetStatus(ctx, id, next, station.Version); err != nil {
		return fmt.Errorf("store station status: %w", err)
	}
	return nil
}

func (s *ObservatoryService) SuspendStation(ctx context.Context, id string) error {
	return s.ChangeStationStatus(ctx, id, model.StationRepair)
}

func (s *ObservatoryService) RestoreStation(ctx context.Context, id string) error {
	return s.ChangeStationStatus(ctx, id, model.StationActive)
}
