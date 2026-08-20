package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/report"
	"github.com/jb843051627/bathypulse/internal/store"
)

func (s *ObservatoryService) SaveCalibration(ctx context.Context, profile model.CalibrationProfile) error {
	if err := profile.Validate(); err != nil {
		return err
	}
	if _, err := s.stations.Get(ctx, profile.StationID); err != nil {
		return fmt.Errorf("calibration station: %w", err)
	}
	calibrations := storeCalibration(s.db)
	if err := calibrations.Disable(ctx, profile.StationID); err != nil {
		return err
	}
	profile.Enabled = true
	return calibrations.Save(ctx, profile)
}

func (s *ObservatoryService) CalibrationReport(ctx context.Context, stationID string, values []float64) (report.CalibrationReport, error) {
	profile, err := storeCalibration(s.db).GetEnabled(ctx, stationID)
	if err != nil {
		return report.CalibrationReport{}, err
	}
	return report.BuildCalibrationReport(*profile, values), nil
}

func storeCalibration(db *store.DB) store.CalibrationRepository {
	return db.Calibrations()
}
