package store

import (
	"context"
	"encoding/json"

	"github.com/jb843051627/bathypulse/internal/model"
)

type CalibrationRepository struct{ db *DB }

func (r CalibrationRepository) Save(ctx context.Context, profile model.CalibrationProfile) error {
	points, err := json.Marshal(profile.Points)
	if err != nil {
		return err
	}
	_, err = r.db.SQL.ExecContext(ctx, `INSERT INTO calibration_profiles(id, station_id, version, points, enabled) VALUES(?, ?, ?, ?, ?)`, profile.ID, profile.StationID, profile.Version, points, profile.Enabled)
	return err
}

func (r CalibrationRepository) GetEnabled(ctx context.Context, stationID string) (*model.CalibrationProfile, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT id, station_id, version, points, enabled FROM calibration_profiles WHERE station_id = ? AND enabled = 1 ORDER BY version DESC LIMIT 1`, stationID)
	var profile model.CalibrationProfile
	var points []byte
	if err := row.Scan(&profile.ID, &profile.StationID, &profile.Version, &points, &profile.Enabled); err != nil {
		return nil, mapNotFound(err)
	}
	if err := json.Unmarshal(points, &profile.Points); err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r CalibrationRepository) Disable(ctx context.Context, stationID string) error {
	_, err := r.db.SQL.ExecContext(ctx, `UPDATE calibration_profiles SET enabled = 0 WHERE station_id = ?`, stationID)
	return err
}
