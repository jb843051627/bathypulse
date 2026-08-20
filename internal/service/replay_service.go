package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/codec"
	"github.com/jb843051627/bathypulse/internal/ingest"
	"github.com/jb843051627/bathypulse/internal/model"
)

func (s *ObservatoryService) ReplayRecord(ctx context.Context, data []byte) error {
	record, err := codec.DecodeReplay(data)
	if err != nil {
		return err
	}
	if record.Kind != "samples" {
		return fmt.Errorf("unsupported replay kind: %w", model.ErrInvalid)
	}
	values, err := codec.DecodeFloat64s(record.Payload)
	if err != nil {
		return err
	}
	return s.SubmitSamples(ctx, ingest.NewBatch(record.StationID, values, record.At))
}

func (s *ObservatoryService) ReplayWindow(ctx context.Context, records []codec.ReplayRecord, start, end time.Time) (int, error) {
	selected := codec.ReplayWindow(records, start, end)
	for _, record := range selected {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		data, err := codec.EncodeReplay(record)
		if err != nil {
			return 0, err
		}
		if err := s.ReplayRecord(ctx, data); err != nil {
			return 0, err
		}
	}
	return len(selected), nil
}

func (s *ObservatoryService) DecodeReplayValues(data []byte) ([]float64, error) {
	return codec.DecodeFloat64s(data)
}
