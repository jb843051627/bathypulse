package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type SnapshotRepository struct{ db *DB }

func (r SnapshotRepository) Save(ctx context.Context, stationID string, at time.Time, health model.StationHealth) error {
	data, err := json.Marshal(health)
	if err != nil {
		return err
	}
	_, err = r.db.SQL.ExecContext(ctx, `INSERT INTO station_snapshots(station_id, observed_at, payload) VALUES(?, ?, ?)`, stationID, formatTime(at), data)
	return err
}

func (r SnapshotRepository) Latest(ctx context.Context, stationID string) (*model.StationHealth, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT payload FROM station_snapshots WHERE station_id = ? ORDER BY observed_at DESC LIMIT 1`, stationID)
	var data []byte
	if err := row.Scan(&data); err != nil {
		return nil, mapNotFound(err)
	}
	var health model.StationHealth
	if err := json.Unmarshal(data, &health); err != nil {
		return nil, err
	}
	return &health, nil
}

func (r SnapshotRepository) RemoveBefore(ctx context.Context, stationID string, before time.Time) (int64, error) {
	result, err := r.db.SQL.ExecContext(ctx, `DELETE FROM station_snapshots WHERE station_id = ? AND observed_at < ?`, stationID, formatTime(before))
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
