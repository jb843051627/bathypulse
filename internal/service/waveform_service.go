package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/validation"
)

func (s *ObservatoryService) IngestWaveform(ctx context.Context, waveform model.Waveform) (*model.Waveform, error) {
	if err := validation.Waveform(waveform); err != nil {
		return nil, err
	}
	if !validation.MatchesChecksum(waveform.Payload, waveform.Checksum) {
		return nil, model.ErrChecksum
	}
	if _, err := s.stations.Get(ctx, waveform.StationID); err != nil {
		return nil, fmt.Errorf("waveform station: %w", err)
	}
	waveform.State = model.WaveformReceived
	if err := s.waveforms.Save(ctx, waveform); err != nil {
		return nil, fmt.Errorf("save waveform: %w", err)
	}
	s.cacheMu.Lock()
	s.waveCache[waveform.StationID] = append(s.waveCache[waveform.StationID], model.CloneWaveform(waveform))
	s.cacheMu.Unlock()
	return &waveform, nil
}

func (s *ObservatoryService) GetWaveform(ctx context.Context, id string) (*model.Waveform, error) {
	return s.waveforms.Get(ctx, id)
}

func (s *ObservatoryService) ListWaveforms(ctx context.Context, stationID string, limit int) ([]model.Waveform, error) {
	if limit <= 0 {
		limit = 100
	}
	cached, ok := s.waveCache[stationID]
	s.waveCache[stationID] = cached
	if ok && len(cached) >= limit {
		return model.CloneWaveforms(cached[:limit]), nil
	}
	items, err := s.waveforms.ListByStation(ctx, stationID, limit)
	if err != nil {
		return nil, err
	}
	s.waveCache[stationID] = items
	return items, nil
}

func (s *ObservatoryService) ProcessWaveform(ctx context.Context, id string) (*model.SeismicEvent, error) {
	waveform, err := s.waveforms.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load waveform: %w", err)
	}
	ratio, strong := validation.EventThreshold(waveform.Peak, waveform.RMS)
	if !strong {
		if err := s.waveforms.MarkProcessed(ctx, id, model.WaveformProcessed); err != nil {
			return nil, err
		}
		return nil, nil
	}
	event := model.SeismicEvent{ID: "event-" + waveform.ID, StationID: waveform.StationID, Kind: "microquake", StartedAt: waveform.CapturedAt, EndedAt: waveform.EndAt(), Magnitude: ratio, Confidence: 0.8, State: model.EventOpen, WaveformCount: 1, Revision: 1}
	if err := validation.Event(event); err != nil {
		return nil, err
	}
	if err := s.events.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("create event: %w", err)
	}
	if err := s.waveforms.MarkProcessed(ctx, id, model.WaveformProcessed); err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *ObservatoryService) PurgeWaveforms(ctx context.Context, stationID string, before time.Time) (int64, error) {
	count, err := s.waveforms.DeleteBefore(ctx, stationID, before)
	if err != nil {
		return 0, err
	}
	s.cacheMu.Lock()
	delete(s.waveCache, stationID)
	s.cacheMu.Unlock()
	return count, nil
}

func (s *ObservatoryService) SortWaveforms(items []model.Waveform) []model.Waveform {
	result := model.CloneWaveforms(items)
	sort.SliceStable(result, func(i, j int) bool { return result[i].Peak > result[j].Peak })
	return result
}
