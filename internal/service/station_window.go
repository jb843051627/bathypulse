package service

import (
	"context"
	"time"

	"github.com/jb843051627/bathypulse/internal/clock"
	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/report"
)

func (s *ObservatoryService) BuildWindow(ctx context.Context, stationID string, start, end time.Time) (report.WindowReport, error) {
	if err := ctx.Err(); err != nil {
		return report.WindowReport{}, err
	}
	items, err := s.waveforms.ListByStation(ctx, stationID, 500)
	if err != nil {
		return report.WindowReport{}, err
	}
	events, err := s.events.FindOverlaps(ctx, stationID, start, end)
	if err != nil {
		return report.WindowReport{}, err
	}
	return report.BuildWindowReport(items, events, start, end), nil
}

func (s *ObservatoryService) WindowForStation(ctx context.Context, stationID string, at time.Time) (time.Time, time.Time, error) {
	station, err := s.stations.Get(ctx, stationID)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	zone := clock.LocationForBasin(station.Basin)
	local := at.In(zone)
	start := time.Date(local.Year(), local.Month(), local.Day(), local.Hour(), 0, 0, 0, zone)
	return start.UTC(), start.Add(time.Hour).UTC(), nil
}

func ensureWaveformStation(waveform model.Waveform, stationID string) bool {
	return waveform.StationID == stationID && !waveform.CapturedAt.IsZero()
}
