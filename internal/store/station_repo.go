package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type StationRepository struct{ db *DB }

func (r StationRepository) Create(ctx context.Context, station model.Station) error {
	_, err := r.db.SQL.ExecContext(ctx, `INSERT INTO stations(id, code, basin, latitude, longitude, depth_meters, status, last_seen, version) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`, station.ID, station.Code, station.Basin, station.Latitude, station.Longitude, station.DepthMeters, station.Status, formatTime(station.LastSeen), station.Version)
	return err
}

func (r StationRepository) Get(ctx context.Context, id string) (*model.Station, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT id, code, basin, latitude, longitude, depth_meters, status, last_seen, version FROM stations WHERE id = ?`, id)
	var station model.Station
	var status, seen string
	err := row.Scan(&station.ID, &station.Code, &station.Basin, &station.Latitude, &station.Longitude, &station.DepthMeters, &status, &seen, &station.Version)
	if err != nil {
		return nil, mapNotFound(err)
	}
	station.Status = model.StationStatus(status)
	station.LastSeen = parseTime(seen)
	return &station, nil
}

func (r StationRepository) FindByCode(ctx context.Context, code string) (*model.Station, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT id, code, basin, latitude, longitude, depth_meters, status, last_seen, version FROM stations WHERE code = ?`, code)
	var station model.Station
	var status, seen string
	err := row.Scan(&station.ID, &station.Code, &station.Basin, &station.Latitude, &station.Longitude, &station.DepthMeters, &status, &seen, &station.Version)
	if err != nil {
		return nil, mapNotFound(err)
	}
	station.Status = model.StationStatus(status)
	station.LastSeen = parseTime(seen)
	return &station, nil
}

func (r StationRepository) List(ctx context.Context, status model.StationStatus) ([]model.Station, error) {
	query := `SELECT id, code, basin, latitude, longitude, depth_meters, status, last_seen, version FROM stations ORDER BY code`
	args := []any{}
	if status != "" {
		query = `SELECT id, code, basin, latitude, longitude, depth_meters, status, last_seen, version FROM stations WHERE status = ? ORDER BY code`
		args = append(args, status)
	}
	rows, err := r.db.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.Station, 0)
	for rows.Next() {
		var station model.Station
		var state, seen string
		if err := rows.Scan(&station.ID, &station.Code, &station.Basin, &station.Latitude, &station.Longitude, &station.DepthMeters, &state, &seen, &station.Version); err != nil {
			return nil, err
		}
		station.Status = model.StationStatus(state)
		station.LastSeen = parseTime(seen)
		result = append(result, station)
	}
	return result, rows.Err()
}

func (r StationRepository) UpdateHeartbeat(ctx context.Context, id string, seen time.Time, version int) error {
	result, err := r.db.SQL.ExecContext(ctx, `UPDATE stations SET last_seen = ?, version = version + 1 WHERE id = ? AND version = ?`, formatTime(seen), id, version)
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

func (r StationRepository) SetStatus(ctx context.Context, id string, next model.StationStatus, version int) error {
	result, err := r.db.SQL.ExecContext(ctx, `UPDATE stations SET status = ?, version = version + 1 WHERE id = ? AND version = ?`, next, id, version)
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

func (r StationRepository) CountByStatus(ctx context.Context, status model.StationStatus) (int, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT COUNT(*) FROM stations WHERE status = ?`, status)
	var count int
	return count, row.Scan(&count)
}

var _ = sql.ErrNoRows
