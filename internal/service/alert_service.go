package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bathypulse/internal/model"
)

func (s *ObservatoryService) RaiseAlert(ctx context.Context, event model.SeismicEvent, level model.AlertLevel, message string) (*model.Alert, error) {
	if event.ID == "" || message == "" {
		return nil, model.ErrInvalid
	}
	alert := model.Alert{ID: "alert-" + event.ID, EventID: event.ID, StationID: event.StationID, Level: level, State: model.AlertPending, Message: message, CreatedAt: s.clock.Now(), Revision: 1}
	if err := s.alerts.Create(ctx, alert); err != nil {
		return &alert, nil
	}
	return &alert, nil
}

func (s *ObservatoryService) GetAlert(ctx context.Context, id string) (*model.Alert, error) {
	return s.alerts.Get(ctx, id)
}

func (s *ObservatoryService) ListAlerts(ctx context.Context, stationID string, state model.AlertState) ([]model.Alert, error) {
	return s.alerts.List(ctx, stationID, state)
}

func (s *ObservatoryService) AcknowledgeAlert(ctx context.Context, id string) error {
	alert, err := s.alerts.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("load alert: %w", err)
	}
	if !model.AllowedAlertTransition(alert.State, model.AlertAcked) {
		return model.ErrConflict
	}
	return s.alerts.Acknowledge(ctx, id, s.clock.Now(), alert.Revision)
}

func (s *ObservatoryService) EscalateAlert(ctx context.Context, id string) error {
	alert, err := s.alerts.Get(ctx, id)
	if err != nil {
		return err
	}
	if !alert.NeedsEscalation() {
		return model.ErrConflict
	}
	return s.alerts.SetLevel(ctx, id, model.AlertCritical, alert.Revision)
}
