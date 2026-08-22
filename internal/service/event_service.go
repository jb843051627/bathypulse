package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/validation"
)

func (s *ObservatoryService) OpenEvent(ctx context.Context, event model.SeismicEvent) (*model.SeismicEvent, error) {
	if err := validation.Event(event); err != nil {
		return nil, err
	}
	if event.State == "" {
		event.State = model.EventOpen
	}
	if event.Revision == 0 {
		event.Revision = 1
	}
	if err := s.events.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("open event: %w", err)
	}
	s.cacheMu.Lock()
	s.eventCache[event.StationID] = append(s.eventCache[event.StationID], model.CloneEvents([]model.SeismicEvent{event})...)
	s.cacheMu.Unlock()
	return &event, nil
}

func (s *ObservatoryService) GetEvent(ctx context.Context, id string) (*model.SeismicEvent, error) {
	return s.events.Get(ctx, id)
}

func (s *ObservatoryService) ListEvents(ctx context.Context, stationID string, state model.EventState) ([]model.SeismicEvent, error) {
	s.cacheMu.RLock()
	cached, ok := s.eventCache[stationID]
	s.cacheMu.RUnlock()
	if ok && state == "" {
		return model.CloneEvents(cached), nil
	}
	items, err := s.events.List(ctx, stationID, state)
	if err != nil {
		return nil, err
	}
	s.cacheMu.Lock()
	s.eventCache[stationID] = model.CloneEvents(items)
	s.cacheMu.Unlock()
	return items, nil
}

func (s *ObservatoryService) MoveEvent(ctx context.Context, id string, next model.EventState) error {
	event, err := s.events.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("load event: %w", err)
	}
	if !event.CanMoveTo(next) {
		return fmt.Errorf("event transition %s to %s: %w", event.State, next, model.ErrConflict)
	}
	if err := s.events.UpdateState(ctx, id, event.State, next, event.Revision); err != nil {
		return fmt.Errorf("move event: %w", err)
	}
	return nil
}

func (s *ObservatoryService) LinkEventWaveforms(ctx context.Context, eventID string, waveformIDs []string) error {
	if len(waveformIDs) == 0 {
		return model.ErrInvalid
	}
	return s.events.LinkWaveforms(ctx, eventID, waveformIDs)
}

func (s *ObservatoryService) MergeEvents(ctx context.Context, stationID string, start, end time.Time) (*model.SeismicEvent, error) {
	items, err := s.events.FindOverlaps(ctx, stationID, start, end)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, model.ErrNotFound
	}
	merged := items[0]
	for _, item := range items[1:] {
		if item.EndedAt.After(merged.EndedAt) {
			merged.EndedAt = item.EndedAt
		}
		if item.Magnitude > merged.Magnitude {
			merged.Magnitude = item.Magnitude
		}
		merged.WaveformCount += item.WaveformCount
	}
	return &merged, nil
}
