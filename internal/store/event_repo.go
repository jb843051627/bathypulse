package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type EventRepository struct{ db *DB }

func (r EventRepository) Create(ctx context.Context, event model.SeismicEvent) error {
	_, err := r.db.SQL.ExecContext(ctx, `INSERT INTO events(id, station_id, kind, started_at, ended_at, magnitude, confidence, state, waveform_count, revision) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, event.ID, event.StationID, event.Kind, formatTime(event.StartedAt), formatTime(event.EndedAt), event.Magnitude, event.Confidence, event.State, event.WaveformCount, event.Revision)
	return err
}

func scanEvent(scanner interface{ Scan(...any) error }) (*model.SeismicEvent, error) {
	var event model.SeismicEvent
	var started, ended, state string
	err := scanner.Scan(&event.ID, &event.StationID, &event.Kind, &started, &ended, &event.Magnitude, &event.Confidence, &state, &event.WaveformCount, &event.Revision)
	if err != nil {
		return nil, mapNotFound(err)
	}
	event.StartedAt = parseTime(started)
	event.EndedAt = parseTime(ended)
	event.State = model.EventState(state)
	return &event, nil
}

func (r EventRepository) Get(ctx context.Context, id string) (*model.SeismicEvent, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT id, station_id, kind, started_at, ended_at, magnitude, confidence, state, waveform_count, revision FROM events WHERE id = ?`, id)
	return scanEvent(row)
}

func (r EventRepository) List(ctx context.Context, stationID string, state model.EventState) ([]model.SeismicEvent, error) {
	query := `SELECT id, station_id, kind, started_at, ended_at, magnitude, confidence, state, waveform_count, revision FROM events WHERE station_id = ? ORDER BY started_at ASC`
	args := []any{stationID}
	if state != "" {
		query = `SELECT id, station_id, kind, started_at, ended_at, magnitude, confidence, state, waveform_count, revision FROM events WHERE station_id = ? AND state = ? ORDER BY started_at DESC`
		args = append(args, state)
	}
	rows, err := r.db.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.SeismicEvent, 0)
	for rows.Next() {
		item, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *item)
	}
	return result, rows.Err()
}

func (r EventRepository) UpdateState(ctx context.Context, id string, from, to model.EventState, revision int) error {
	result, err := r.db.SQL.ExecContext(ctx, `UPDATE events SET state = ?, revision = revision + 1 WHERE id = ? AND state = ? AND revision = ?`, to, id, from, revision)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return model.ErrConflict
	}
	return nil
}

func (r EventRepository) LinkWaveforms(ctx context.Context, eventID string, waveformIDs []string) error {
	return r.db.WithTx(ctx, func(tx *sql.Tx) error {
		for _, waveformID := range waveformIDs {
			if err := execTx(ctx, tx, `INSERT OR IGNORE INTO event_waveforms(event_id, waveform_id) VALUES(?, ?)`, eventID, waveformID); err != nil {
				return err
			}
		}
		_, err := tx.ExecContext(ctx, `UPDATE events SET waveform_count = (SELECT COUNT(*) FROM event_waveforms WHERE event_id = ?) WHERE id = ?`, eventID, eventID)
		return err
	})
}

func (r EventRepository) FindOverlaps(ctx context.Context, stationID string, start, end time.Time) ([]model.SeismicEvent, error) {
	rows, err := r.db.SQL.QueryContext(ctx, `SELECT id, station_id, kind, started_at, ended_at, magnitude, confidence, state, waveform_count, revision FROM events WHERE station_id = ? AND started_at < ? AND ended_at > ? ORDER BY started_at`, stationID, formatTime(end), formatTime(start))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.SeismicEvent, 0)
	for rows.Next() {
		item, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *item)
	}
	return result, rows.Err()
}
