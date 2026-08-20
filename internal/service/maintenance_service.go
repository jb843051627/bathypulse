package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/validation"
)

func (s *ObservatoryService) PlanMaintenance(ctx context.Context, item model.Maintenance) (*model.Maintenance, error) {
	if err := validation.Maintenance(item); err != nil {
		return nil, err
	}
	if _, err := s.stations.Get(ctx, item.StationID); err != nil {
		return nil, fmt.Errorf("maintenance station: %w", err)
	}
	if item.State == "" {
		item.State = model.MaintenancePlanned
	}
	if item.Revision == 0 {
		item.Revision = 1
	}
	if err := s.maintenance.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("plan maintenance: %w", err)
	}
	return &item, nil
}

func (s *ObservatoryService) GetMaintenance(ctx context.Context, id string) (*model.Maintenance, error) {
	return s.maintenance.Get(ctx, id)
}

func (s *ObservatoryService) StartMaintenance(ctx context.Context, id string) error {
	return s.moveMaintenance(ctx, id, model.MaintenanceRunning)
}

func (s *ObservatoryService) CompleteMaintenance(ctx context.Context, id string) error {
	return s.moveMaintenance(ctx, id, model.MaintenanceFinished)
}

func (s *ObservatoryService) CancelMaintenance(ctx context.Context, id string) error {
	return s.moveMaintenance(ctx, id, model.MaintenanceCanceled)
}

func (s *ObservatoryService) moveMaintenance(ctx context.Context, id string, next model.MaintenanceState) error {
	item, err := s.maintenance.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("load maintenance: %w", err)
	}
	if err := s.maintenance.Transition(ctx, id, item.State, next, item.Revision); err != nil {
		return fmt.Errorf("store maintenance: %w", err)
	}
	return nil
}

func (s *ObservatoryService) DueMaintenance(ctx context.Context) ([]model.Maintenance, error) {
	return s.maintenance.ListDue(ctx, s.clock.Now())
}
