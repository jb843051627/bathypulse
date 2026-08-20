package service

import (
	"context"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

func (s *ObservatoryService) StationHealth(ctx context.Context, id string) (model.StationStatus, error) {
	station, err := s.stations.Get(ctx, id)
	if err != nil {
		return "", err
	}
	if !station.IsUsable(s.clock.Now()) {
		return model.StationOffline, nil
	}
	return station.Status, nil
}

func (s *ObservatoryService) MarkStaleStations(ctx context.Context, age time.Duration) (int, error) {
	stations, err := s.stations.List(ctx, "")
	if err != nil {
		return 0, err
	}
	changed := 0
	cutoff := s.clock.Now().Add(-age)
	for _, station := range stations {
		if station.LastSeen.Before(cutoff) && station.Status == model.StationActive {
			if err := s.stations.SetStatus(ctx, station.ID, model.StationOffline, station.Version); err != nil {
				return changed, err
			}
			changed++
		}
	}
	return changed, nil
}

func (s *ObservatoryService) Metrics() map[string]float64 {
	return s.metrics.Snapshot()
}
