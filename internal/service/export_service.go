package service

import (
	"bytes"
	"context"
	"time"

	"github.com/jb843051627/bathypulse/internal/codec"
)

func (s *ObservatoryService) ExportTimeline(ctx context.Context, stationID string, day time.Time) ([]byte, error) {
	items, err := s.Timeline(ctx, stationID, day)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	if err := codec.WriteTimeline(&buffer, items); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s *ObservatoryService) ExportSamples(ctx context.Context, stationID string, limit int) ([]byte, error) {
	items, err := s.samples.Latest(ctx, stationID, limit)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	for _, item := range items {
		if err := codec.WriteSample(&buffer, item); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}
