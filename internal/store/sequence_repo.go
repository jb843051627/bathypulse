package store

import (
	"context"
	"database/sql"
)

type SequenceRepository struct{ db *DB }

func (r SequenceRepository) Reserve(ctx context.Context, stationID string, count int64) (int64, error) {
	if count <= 0 {
		return 0, nil
	}
	var next int64
	err := r.db.WithTx(ctx, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `SELECT next_sequence FROM station_sequences WHERE station_id = ?`, stationID)
		if err := row.Scan(&next); err != nil {
			if err != sql.ErrNoRows {
				return err
			}
			next = 1
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO station_sequences(station_id, next_sequence) VALUES(?, ?) ON CONFLICT(station_id) DO UPDATE SET next_sequence = excluded.next_sequence`, stationID, next+count); err != nil {
			return err
		}
		return nil
	})
	return next, err
}

func (r SequenceRepository) Current(ctx context.Context, stationID string) (int64, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT next_sequence FROM station_sequences WHERE station_id = ?`, stationID)
	var value int64
	return value, row.Scan(&value)
}
