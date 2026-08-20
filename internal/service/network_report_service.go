package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/report"
)

func (s *ObservatoryService) NetworkReport(ctx context.Context, stationID string) (report.NetworkReport, error) {
	stations, err := s.stations.List(ctx, "")
	if err != nil {
		return report.NetworkReport{}, err
	}
	events, err := s.events.List(ctx, stationID, "")
	if err != nil {
		return report.NetworkReport{}, err
	}
	alerts, err := s.alerts.List(ctx, stationID, "")
	if err != nil {
		return report.NetworkReport{}, err
	}
	health := make(map[string]float64)
	for _, station := range stations {
		snapshot, err := s.LatestHealth(ctx, station.ID)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				continue
			}
			return report.NetworkReport{}, fmt.Errorf("network health: %w", err)
		}
		health[station.ID] = snapshot.Score()
	}
	return report.BuildNetworkReport(stations, events, alerts, health), nil
}

func (s *ObservatoryService) StationSummary(ctx context.Context, stationID string) (model.ObservatorySummary, error) {
	if _, err := s.GetStation(ctx, stationID); err != nil {
		return model.ObservatorySummary{}, err
	}
	return s.Summary(ctx)
}
