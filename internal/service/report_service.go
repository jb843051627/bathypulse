package service

import (
	"context"
	"sort"
	"time"

	"github.com/jb843051627/bathypulse/internal/clock"
	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/report"
)

func (s *ObservatoryService) Summary(ctx context.Context) (model.ObservatorySummary, error) {
	summary, err := s.queries.Summary(ctx)
	if err != nil {
		return summary, err
	}
	summary.GeneratedAt = s.clock.Now()
	return summary, nil
}

func (s *ObservatoryService) Timeline(ctx context.Context, stationID string, day time.Time) ([]model.TimelineItem, error) {
	station, err := s.stations.Get(ctx, stationID)
	if err != nil {
		return nil, err
	}
	start, end := report.Builder{}.DayBounds(day, clock.LocationForBasin(station.Basin))
	events, err := s.events.List(ctx, stationID, "")
	if err != nil {
		return nil, err
	}
	alerts, err := s.alerts.List(ctx, stationID, "")
	if err != nil {
		return nil, err
	}
	items := report.Builder{}.Timeline(events, alerts, station.Basin)
	filtered := make([]model.TimelineItem, 0, len(items))
	for _, item := range items {
		if !item.At.Before(start) && item.At.Before(end) {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (s *ObservatoryService) RankStations(ctx context.Context) ([]model.Station, error) {
	stations, err := s.stations.List(ctx, "")
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, station := range stations {
		_, count, err := s.samples.Aggregate(ctx, station.ID, time.Time{}, s.clock.Now().Add(time.Minute))
		if err != nil {
			return nil, err
		}
		counts[station.ID] = count
	}
	result := report.RankStations(stations, counts)
	sort.SliceStable(result, func(i, j int) bool { return result[i].Code < result[j].Code })
	return result, nil
}
