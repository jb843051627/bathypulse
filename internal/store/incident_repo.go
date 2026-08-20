package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type IncidentRepository struct{ db *DB }

func (r IncidentRepository) Create(ctx context.Context, item model.Incident) error {
	_, err := r.db.SQL.ExecContext(ctx, `INSERT INTO incidents(id, station_id, alert_id, state, summary, opened_at, closed_at, revision) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`, item.ID, item.StationID, item.AlertID, item.State, item.Summary, formatTime(item.OpenedAt), nil, item.Revision)
	return err
}

func (r IncidentRepository) Get(ctx context.Context, id string) (model.Incident, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT id, station_id, alert_id, state, summary, opened_at, closed_at, revision FROM incidents WHERE id = ?`, id)
	return scanIncident(row)
}

func scanIncident(scanner interface{ Scan(...any) error }) (model.Incident, error) {
	var item model.Incident
	var state, opened string
	var closed sql.NullString
	if err := scanner.Scan(&item.ID, &item.StationID, &item.AlertID, &state, &item.Summary, &opened, &closed, &item.Revision); err != nil {
		return item, mapNotFound(err)
	}
	item.State = model.IncidentState(state)
	item.OpenedAt = parseTime(opened)
	if closed.Valid {
		value := parseTime(closed.String)
		item.ClosedAt = &value
	}
	return item, nil
}

func (r IncidentRepository) Move(ctx context.Context, id string, from, to model.IncidentState, at time.Time, revision int) error {
	var closed any
	if to == model.IncidentClosed {
		closed = formatTime(at)
	}
	result, err := r.db.SQL.ExecContext(ctx, `UPDATE incidents SET state = ?, closed_at = ?, revision = revision + 1 WHERE id = ? AND state = ? AND revision = ?`, to, closed, id, from, revision)
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

func (r IncidentRepository) Open(ctx context.Context, stationID string) ([]model.Incident, error) {
	rows, err := r.db.SQL.QueryContext(ctx, `SELECT id, station_id, alert_id, state, summary, opened_at, closed_at, revision FROM incidents WHERE station_id = ? AND state <> ? ORDER BY opened_at DESC`, stationID, model.IncidentClosed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.Incident, 0)
	for rows.Next() {
		item, err := scanIncident(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
