package store

import (
	"context"
	"time"
)

type RetentionRule struct {
	ID        string
	StationID string
	KeepDays  int
	Enabled   bool
	UpdatedAt time.Time
}

type RetentionRepository struct{ db *DB }

func (r RetentionRepository) Save(ctx context.Context, rule RetentionRule) error {
	_, err := r.db.SQL.ExecContext(ctx, `INSERT INTO retention_rules(id, station_id, keep_days, enabled, updated_at) VALUES(?, ?, ?, ?, ?)`, rule.ID, rule.StationID, rule.KeepDays, rule.Enabled, formatTime(rule.UpdatedAt))
	return err
}

func (r RetentionRepository) Get(ctx context.Context, stationID string) (RetentionRule, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT id, station_id, keep_days, enabled, updated_at FROM retention_rules WHERE station_id = ? AND enabled = 1`, stationID)
	var rule RetentionRule
	var updated string
	if err := row.Scan(&rule.ID, &rule.StationID, &rule.KeepDays, &rule.Enabled, &updated); err != nil {
		return rule, mapNotFound(err)
	}
	rule.UpdatedAt = parseTime(updated)
	return rule, nil
}

func (r RetentionRepository) Disable(ctx context.Context, stationID string) error {
	_, err := r.db.SQL.ExecContext(ctx, `UPDATE retention_rules SET enabled = 0 WHERE station_id = ?`, stationID)
	return err
}

func (r RetentionRepository) RemoveWaveforms(ctx context.Context, stationID string, before time.Time) (int64, error) {
	result, err := r.db.SQL.ExecContext(ctx, `DELETE FROM waveforms WHERE station_id = ? AND captured_at < ?`, stationID, formatTime(before))
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
