package store

import (
	"context"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type QueryRepository struct{ db *DB }

func (r QueryRepository) Summary(ctx context.Context) (model.ObservatorySummary, error) {
	var summary model.ObservatorySummary
	queries := []struct {
		query       string
		destination *int
	}{
		{`SELECT COUNT(*) FROM stations`, &summary.StationCount},
		{`SELECT COUNT(*) FROM stations WHERE status = ?`, &summary.ActiveStations},
		{`SELECT COUNT(*) FROM events WHERE state IN (?, ?)`, &summary.OpenEvents},
		{`SELECT COUNT(*) FROM alerts WHERE state = ?`, &summary.PendingAlerts},
		{`SELECT COUNT(*) FROM samples`, &summary.Samples},
	}
	args := [][]any{nil, {model.StationActive}, {model.EventOpen, model.EventReview}, {model.AlertPending}, nil}
	for i, item := range queries {
		row := r.db.SQL.QueryRowContext(ctx, item.query, args[i]...)
		if err := row.Scan(item.destination); err != nil {
			return summary, err
		}
	}
	summary.GeneratedAt = time.Now().UTC()
	return summary, nil
}

func (r QueryRepository) Timeline(ctx context.Context, stationID string, start, end time.Time) ([]model.TimelineItem, error) {
	items := make([]model.TimelineItem, 0)
	events, err := r.db.Events().List(ctx, stationID, "")
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		if event.StartedAt.Before(start) || !event.StartedAt.Before(end) {
			continue
		}
		items = append(items, model.TimelineItem{At: event.StartedAt, Kind: "event", EntityID: event.ID, Station: stationID, Summary: event.Kind})
	}
	alerts, err := r.db.Alerts().List(ctx, stationID, "")
	if err != nil {
		return nil, err
	}
	for _, alert := range alerts {
		if alert.CreatedAt.Before(start) || !alert.CreatedAt.Before(end) {
			continue
		}
		items = append(items, model.TimelineItem{At: alert.CreatedAt, Kind: "alert", EntityID: alert.ID, Station: stationID, Summary: alert.Message})
	}
	return items, nil
}

func (r QueryRepository) StationActivity(ctx context.Context, stationID string, start, end time.Time) (float64, int, error) {
	return r.db.Samples().Aggregate(ctx, stationID, start, end)
}
