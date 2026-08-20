package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/validation"
)

func (s *ObservatoryService) PlanMaintenanceBatch(ctx context.Context, items []model.Maintenance) error {
	if len(items) == 0 {
		return model.ErrInvalid
	}
	return s.db.WithTx(ctx, func(tx *sql.Tx) error {
		for _, item := range items {
			if err := validation.Maintenance(item); err != nil {
				return err
			}
			if _, err := s.stations.Get(ctx, item.StationID); err != nil {
				return fmt.Errorf("batch maintenance station: %w", err)
			}
			if item.State == "" {
				item.State = model.MaintenancePlanned
			}
			if item.Revision == 0 {
				item.Revision = 1
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO maintenance(id, station_id, window_start, window_end, state, reason, revision) VALUES(?, ?, ?, ?, ?, ?, ?)`, item.ID, item.StationID, item.WindowStart.UTC().Format(time.RFC3339Nano), item.WindowEnd.UTC().Format(time.RFC3339Nano), item.State, item.Reason, item.Revision); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *ObservatoryService) CancelDueMaintenance(ctx context.Context) (int, error) {
	items, err := s.DueMaintenance(ctx)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, item := range items {
		if err := s.CancelMaintenance(ctx, item.ID); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *ObservatoryService) MaintenanceReady(ctx context.Context, id string) (bool, error) {
	item, err := s.GetMaintenance(ctx, id)
	if err != nil {
		return false, err
	}
	station, err := s.GetStation(ctx, item.StationID)
	if err != nil {
		return false, err
	}
	return station.Status != model.StationOffline && station.Status != model.StationRepair, nil
}
