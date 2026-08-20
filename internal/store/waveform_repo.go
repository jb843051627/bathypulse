package store

import (
	"context"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type WaveformRepository struct{ db *DB }

func (r WaveformRepository) Save(ctx context.Context, waveform model.Waveform) error {
	_, err := r.db.SQL.ExecContext(ctx, `INSERT INTO waveforms(id, station_id, captured_at, duration_ms, sample_rate, peak, rms, payload, checksum, state, sequence) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, waveform.ID, waveform.StationID, formatTime(waveform.CapturedAt), waveform.DurationMs, waveform.SampleRate, waveform.Peak, waveform.RMS, waveform.Payload, waveform.Checksum, waveform.State, waveform.Sequence)
	return err
}

func (r WaveformRepository) Get(ctx context.Context, id string) (*model.Waveform, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT id, station_id, captured_at, duration_ms, sample_rate, peak, rms, payload, checksum, state, sequence FROM waveforms WHERE id = ?`, id)
	var item model.Waveform
	var captured, state string
	err := row.Scan(&item.ID, &item.StationID, &captured, &item.DurationMs, &item.SampleRate, &item.Peak, &item.RMS, &item.Payload, &item.Checksum, &state, &item.Sequence)
	if err != nil {
		return nil, mapNotFound(err)
	}
	item.CapturedAt = parseTime(captured)
	item.State = model.WaveformState(state)
	item.Payload = model.CloneBytes(item.Payload)
	return &item, nil
}

func (r WaveformRepository) ListByStation(ctx context.Context, stationID string, limit int) ([]model.Waveform, error) {
	rows, err := r.db.SQL.QueryContext(ctx, `SELECT id, station_id, captured_at, duration_ms, sample_rate, peak, rms, payload, checksum, state, sequence FROM waveforms WHERE station_id = ? ORDER BY captured_at ASC LIMIT ?`, stationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.Waveform, 0)
	for rows.Next() {
		var item model.Waveform
		var captured, state string
		if err := rows.Scan(&item.ID, &item.StationID, &captured, &item.DurationMs, &item.SampleRate, &item.Peak, &item.RMS, &item.Payload, &item.Checksum, &state, &item.Sequence); err != nil {
			return nil, err
		}
		item.CapturedAt = parseTime(captured)
		item.State = model.WaveformState(state)
		item.Payload = model.CloneBytes(item.Payload)
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r WaveformRepository) MarkProcessed(ctx context.Context, id string, state model.WaveformState) error {
	_, err := r.db.SQL.ExecContext(ctx, `UPDATE waveforms SET state = ? WHERE id = ?`, state, id)
	return err
}

func (r WaveformRepository) DeleteBefore(ctx context.Context, stationID string, before time.Time) (int64, error) {
	result, err := r.db.SQL.ExecContext(ctx, `DELETE FROM waveforms WHERE station_id = ? AND captured_at < ?`, stationID, formatTime(before))
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r WaveformRepository) Count(ctx context.Context, stationID string) (int, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT COUNT(*) FROM waveforms WHERE station_id = ?`, stationID)
	var count int
	return count, row.Scan(&count)
}
