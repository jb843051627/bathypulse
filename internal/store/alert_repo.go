package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type AlertRepository struct{ db *DB }

func (r AlertRepository) Create(ctx context.Context, alert model.Alert) error {
	_, err := r.db.SQL.ExecContext(ctx, `INSERT OR IGNORE INTO alerts(id, event_id, station_id, level, state, message, created_at, acked_at, revision) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`, alert.ID, alert.EventID, alert.StationID, alert.Level, alert.State, alert.Message, formatTime(alert.CreatedAt), nil, alert.Revision)
	return err
}

func (r AlertRepository) Get(ctx context.Context, id string) (*model.Alert, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT id, event_id, station_id, level, state, message, created_at, acked_at, revision FROM alerts WHERE id = ?`, id)
	return scanAlert(row)
}

func scanAlert(scanner interface{ Scan(...any) error }) (*model.Alert, error) {
	var alert model.Alert
	var level, state, created, acked sql.NullString
	err := scanner.Scan(&alert.ID, &alert.EventID, &alert.StationID, &level, &state, &alert.Message, &created, &acked, &alert.Revision)
	if err != nil {
		return nil, mapNotFound(err)
	}
	alert.Level = model.AlertLevel(level.String)
	alert.State = model.AlertState(state.String)
	alert.CreatedAt = parseTime(created.String)
	if acked.Valid {
		value := parseTime(acked.String)
		alert.AckedAt = &value
	}
	return &alert, nil
}

func (r AlertRepository) List(ctx context.Context, stationID string, state model.AlertState) ([]model.Alert, error) {
	rows, err := r.db.SQL.QueryContext(ctx, `SELECT id, event_id, station_id, level, state, message, created_at, acked_at, revision FROM alerts WHERE station_id = ? AND (? = '' OR state = ?) ORDER BY created_at DESC`, stationID, state, state)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.Alert, 0)
	for rows.Next() {
		item, err := scanAlert(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *item)
	}
	return result, rows.Err()
}

func (r AlertRepository) Acknowledge(ctx context.Context, id string, at time.Time, revision int) error {
	result, err := r.db.SQL.ExecContext(ctx, `UPDATE alerts SET state = ?, acked_at = ?, revision = revision + 1 WHERE id = ? AND state = ? AND revision = ?`, model.AlertAcked, formatTime(at), id, model.AlertPending, revision)
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

func (r AlertRepository) SetLevel(ctx context.Context, id string, level model.AlertLevel, revision int) error {
	result, err := r.db.SQL.ExecContext(ctx, `UPDATE alerts SET level = ?, revision = revision + 1 WHERE id = ? AND revision = ?`, level, id, revision)
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

func (r AlertRepository) CountPending(ctx context.Context) (int, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT COUNT(*) FROM alerts WHERE state = ?`, model.AlertPending)
	var count int
	return count, row.Scan(&count)
}
