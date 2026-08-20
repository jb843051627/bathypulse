package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/report"
	"github.com/jb843051627/bathypulse/internal/validation"
)

func (s *ObservatoryService) OpenIncident(ctx context.Context, alertID string, summary string) (model.Incident, error) {
	alert, err := s.alerts.Get(ctx, alertID)
	if err != nil {
		return model.Incident{}, fmt.Errorf("incident alert: %w", err)
	}
	item := model.Incident{ID: "incident-" + alert.ID, StationID: alert.StationID, AlertID: alert.ID, State: model.IncidentOpen, Summary: summary, OpenedAt: s.clock.Now(), Revision: 1}
	if err := validation.Incident(item); err != nil {
		return model.Incident{}, err
	}
	if err := s.db.Incidents().Create(ctx, item); err != nil {
		return model.Incident{}, fmt.Errorf("open incident: %w", err)
	}
	return item, nil
}

func (s *ObservatoryService) GetIncident(ctx context.Context, id string) (model.Incident, error) {
	return s.db.Incidents().Get(ctx, id)
}

func (s *ObservatoryService) InvestigateIncident(ctx context.Context, id string) error {
	return s.moveIncident(ctx, id, model.IncidentInvestig)
}

func (s *ObservatoryService) CloseIncident(ctx context.Context, id string) error {
	return s.moveIncident(ctx, id, model.IncidentClosed)
}

func (s *ObservatoryService) moveIncident(ctx context.Context, id string, next model.IncidentState) error {
	item, err := s.db.Incidents().Get(ctx, id)
	if err != nil {
		return err
	}
	if !item.CanMoveTo(next) {
		return model.ErrConflict
	}
	return s.db.Incidents().Move(ctx, id, item.State, next, s.clock.Now(), item.Revision)
}

func (s *ObservatoryService) IncidentReport(ctx context.Context, stationID string) (report.IncidentReport, error) {
	incidents, err := s.db.Incidents().Open(ctx, stationID)
	if err != nil {
		return report.IncidentReport{}, err
	}
	alerts, err := s.alerts.List(ctx, stationID, "")
	if err != nil {
		return report.IncidentReport{}, err
	}
	return report.BuildIncidentReport(stationID, incidents, alerts), nil
}

func (s *ObservatoryService) IncidentCount(ctx context.Context, stationID string) (int, error) {
	items, err := s.db.Incidents().Open(ctx, stationID)
	return len(items), err
}
