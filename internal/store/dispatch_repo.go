package store

import (
	"context"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type DispatchRepository struct{ db *DB }

func (r DispatchRepository) Create(ctx context.Context, job model.DispatchJob) error {
	_, err := r.db.SQL.ExecContext(ctx, `INSERT INTO dispatch_jobs(id, station_id, kind, state, attempt, created_at, updated_at, payload) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`, job.ID, job.StationID, job.Kind, job.State, job.Attempt, formatTime(job.CreatedAt), formatTime(job.UpdatedAt), job.Payload)
	return err
}

func (r DispatchRepository) Get(ctx context.Context, id string) (model.DispatchJob, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT id, station_id, kind, state, attempt, created_at, updated_at, payload FROM dispatch_jobs WHERE id = ?`, id)
	var job model.DispatchJob
	var state, created, updated string
	if err := row.Scan(&job.ID, &job.StationID, &job.Kind, &state, &job.Attempt, &created, &updated, &job.Payload); err != nil {
		return job, mapNotFound(err)
	}
	job.State = model.DispatchState(state)
	job.CreatedAt = parseTime(created)
	job.UpdatedAt = parseTime(updated)
	return job, nil
}

func (r DispatchRepository) Move(ctx context.Context, id string, from, to model.DispatchState, at time.Time, attempt int) error {
	result, err := r.db.SQL.ExecContext(ctx, `UPDATE dispatch_jobs SET state = ?, updated_at = ?, attempt = ? WHERE id = ? AND state = ?`, to, formatTime(at), attempt, id, from)
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

func (r DispatchRepository) Queued(ctx context.Context, limit int) ([]model.DispatchJob, error) {
	rows, err := r.db.SQL.QueryContext(ctx, `SELECT id, station_id, kind, state, attempt, created_at, updated_at, payload FROM dispatch_jobs WHERE state = ? ORDER BY created_at LIMIT ?`, model.DispatchQueued, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.DispatchJob, 0)
	for rows.Next() {
		var job model.DispatchJob
		var state, created, updated string
		if err := rows.Scan(&job.ID, &job.StationID, &job.Kind, &state, &job.Attempt, &created, &updated, &job.Payload); err != nil {
			return nil, err
		}
		job.State = model.DispatchState(state)
		job.CreatedAt = parseTime(created)
		job.UpdatedAt = parseTime(updated)
		result = append(result, job)
	}
	return result, rows.Err()
}
