package store

import (
	"context"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type MaintenanceRepository struct{ db *DB }

func (r MaintenanceRepository) Create(ctx context.Context, item model.Maintenance) error {
	_, err := r.db.SQL.ExecContext(ctx, `INSERT INTO maintenance(id, station_id, window_start, window_end, state, reason, revision) VALUES(?, ?, ?, ?, ?, ?, ?)`, item.ID, item.StationID, formatTime(item.WindowStart), formatTime(item.WindowEnd), item.State, item.Reason, item.Revision)
	return err
}

func (r MaintenanceRepository) Get(ctx context.Context, id string) (*model.Maintenance, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT id, station_id, window_start, window_end, state, reason, revision FROM maintenance WHERE id = ?`, id)
	return scanMaintenance(row)
}

func scanMaintenance(scanner interface{ Scan(...any) error }) (*model.Maintenance, error) {
	var item model.Maintenance
	var start, end, state string
	err := scanner.Scan(&item.ID, &item.StationID, &start, &end, &state, &item.Reason, &item.Revision)
	if err != nil {
		return nil, mapNotFound(err)
	}
	item.WindowStart = parseTime(start)
	item.WindowEnd = parseTime(end)
	item.State = model.MaintenanceState(state)
	return &item, nil
}

func (r MaintenanceRepository) ListDue(ctx context.Context, before time.Time) ([]model.Maintenance, error) {
	rows, err := r.db.SQL.QueryContext(ctx, `SELECT id, station_id, window_start, window_end, state, reason, revision FROM maintenance WHERE window_start <= ? AND state IN (?, ?) ORDER BY window_start`, formatTime(before), model.MaintenancePlanned, model.MaintenanceRunning)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.Maintenance, 0)
	for rows.Next() {
		item, err := scanMaintenance(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *item)
	}
	return result, rows.Err()
}

func (r MaintenanceRepository) Transition(ctx context.Context, id string, from, to model.MaintenanceState, revision int) error {
	result, err := r.db.SQL.ExecContext(ctx, `UPDATE maintenance SET state = ?, revision = revision + 1 WHERE id = ? AND state = ? AND revision = ?`, to, id, from, revision)
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
