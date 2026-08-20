package store

import (
	"context"
	"time"

	"github.com/jb843051627/bathypulse/internal/codec"
)

type ReplayRepository struct{ db *DB }

func (r ReplayRepository) Save(ctx context.Context, id string, record codec.ReplayRecord) error {
	data, err := codec.EncodeReplay(record)
	if err != nil {
		return err
	}
	_, err = r.db.SQL.ExecContext(ctx, `INSERT INTO replay_records(id, station_id, kind, occurred_at, payload) VALUES(?, ?, ?, ?, ?)`, id, record.StationID, record.Kind, formatTime(record.At), data)
	return err
}

func (r ReplayRepository) List(ctx context.Context, stationID string, start, end time.Time) ([]codec.ReplayRecord, error) {
	rows, err := r.db.SQL.QueryContext(ctx, `SELECT payload FROM replay_records WHERE station_id = ? AND occurred_at >= ? AND occurred_at < ? ORDER BY occurred_at`, stationID, formatTime(start), formatTime(end))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]codec.ReplayRecord, 0)
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}
		record, err := codec.DecodeReplay(data)
		if err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return result, rows.Err()
}

func (r ReplayRepository) DeleteBefore(ctx context.Context, stationID string, before time.Time) (int64, error) {
	result, err := r.db.SQL.ExecContext(ctx, `DELETE FROM replay_records WHERE station_id = ? AND occurred_at < ?`, stationID, formatTime(before))
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
