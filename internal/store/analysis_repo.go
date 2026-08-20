package store

import (
	"context"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type AnalysisRepository struct{ db *DB }

func (r AnalysisRepository) Create(ctx context.Context, run model.AnalysisRun) error {
	_, err := r.db.SQL.ExecContext(ctx, `INSERT INTO analysis_runs(id, station_id, started_at, finished_at, state, samples, events, score) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`, run.ID, run.StationID, formatTime(run.StartedAt), formatTime(run.FinishedAt), run.State, run.Samples, run.Events, run.Score)
	return err
}

func (r AnalysisRepository) Get(ctx context.Context, id string) (model.AnalysisRun, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT id, station_id, started_at, finished_at, state, samples, events, score FROM analysis_runs WHERE id = ?`, id)
	var run model.AnalysisRun
	var started, finished string
	if err := row.Scan(&run.ID, &run.StationID, &started, &finished, &run.State, &run.Samples, &run.Events, &run.Score); err != nil {
		return run, mapNotFound(err)
	}
	run.StartedAt = parseTime(started)
	run.FinishedAt = parseTime(finished)
	return run, nil
}

func (r AnalysisRepository) Complete(ctx context.Context, id string, finished time.Time, score float64) error {
	_, err := r.db.SQL.ExecContext(ctx, `UPDATE analysis_runs SET finished_at = ?, state = ?, score = ? WHERE id = ? AND state = ?`, formatTime(finished), "complete", score, id, "running")
	return err
}

func (r AnalysisRepository) List(ctx context.Context, stationID string, limit int) ([]model.AnalysisRun, error) {
	rows, err := r.db.SQL.QueryContext(ctx, `SELECT id, station_id, started_at, finished_at, state, samples, events, score FROM analysis_runs WHERE station_id = ? ORDER BY started_at DESC LIMIT ?`, stationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.AnalysisRun, 0)
	for rows.Next() {
		var run model.AnalysisRun
		var started, finished string
		if err := rows.Scan(&run.ID, &run.StationID, &started, &finished, &run.State, &run.Samples, &run.Events, &run.Score); err != nil {
			return nil, err
		}
		run.StartedAt = parseTime(started)
		run.FinishedAt = parseTime(finished)
		result = append(result, run)
	}
	return result, rows.Err()
}
