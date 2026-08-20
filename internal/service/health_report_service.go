package service

import (
	"context"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/report"
)

func (s *ObservatoryService) BuildHealthReport(ctx context.Context) (report.HealthReport, error) {
	stations, err := s.stations.List(ctx, "")
	if err != nil {
		return report.HealthReport{}, err
	}
	items := make([]model.StationHealth, 0, len(stations))
	for _, station := range stations {
		health, err := s.LatestHealth(ctx, station.ID)
		if err != nil {
			if err == model.ErrNotFound {
				continue
			}
			return report.HealthReport{}, err
		}
		items = append(items, *health)
	}
	return report.BuildHealthReport(s.clock.Now(), items), nil
}

func (s *ObservatoryService) HealthGrade(ctx context.Context) (string, error) {
	result, err := s.BuildHealthReport(ctx)
	if err != nil {
		return "", err
	}
	return report.HealthGrade(result.Average), nil
}

func (s *ObservatoryService) HealthReady(ctx context.Context) (bool, error) {
	result, err := s.BuildHealthReport(ctx)
	if err != nil {
		return false, err
	}
	return result.Ready == len(result.Stations) && len(result.Stations) > 0, nil
}
