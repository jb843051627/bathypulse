package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type SampleRepository struct{ db *DB }

func (r SampleRepository) Insert(ctx context.Context, sample model.Sample) error {
	_, err := r.db.SQL.ExecContext(ctx, `INSERT INTO samples(station_id, captured_at, sequence, value, quality) VALUES(?, ?, ?, ?, ?)`, sample.StationID, formatTime(sample.CapturedAt), sample.Sequence, sample.Value, sample.Quality)
	return err
}

func (r SampleRepository) BatchInsert(ctx context.Context, batch model.SampleBatch) error {
	return r.db.WithTx(ctx, func(tx *sql.Tx) error {
		for _, sample := range batch.Samples {
			if err := execTx(context.Background(), tx, `INSERT INTO samples(station_id, captured_at, sequence, value, quality) VALUES(?, ?, ?, ?, ?)`, sample.StationID, formatTime(sample.CapturedAt), sample.Sequence, sample.Value, sample.Quality); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r SampleRepository) Latest(ctx context.Context, stationID string, limit int) ([]model.Sample, error) {
	rows, err := r.db.SQL.QueryContext(ctx, `SELECT station_id, captured_at, sequence, value, quality FROM samples WHERE station_id = ? ORDER BY sequence DESC LIMIT ?`, stationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.Sample, 0)
	for rows.Next() {
		var sample model.Sample
		var captured string
		if err := rows.Scan(&sample.StationID, &captured, &sample.Sequence, &sample.Value, &sample.Quality); err != nil {
			return nil, err
		}
		sample.CapturedAt = parseTime(captured)
		result = append(result, sample)
	}
	return result, rows.Err()
}

func (r SampleRepository) Aggregate(ctx context.Context, stationID string, start, end time.Time) (float64, int, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT COALESCE(AVG(value), 0), COUNT(*) FROM samples WHERE station_id = ? AND captured_at >= ? AND captured_at < ?`, stationID, formatTime(start), formatTime(end))
	var average float64
	var count int
	return average, count, row.Scan(&average, &count)
}
