package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bathypulse/internal/ingest"
	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/validation"
)

func (s *ObservatoryService) SubmitSamples(ctx context.Context, batch model.SampleBatch) error {
	if err := validation.Batch(batch); err != nil {
		return err
	}
	s.cacheMu.RLock()
	processor := s.processor
	s.cacheMu.RUnlock()
	if processor == nil {
		return s.persistSamples(ctx, batch)
	}
	if err := processor.Submit(ctx, batch); err != nil {
		return fmt.Errorf("submit samples: %w", err)
	}
	return nil
}

func (s *ObservatoryService) persistSamples(ctx context.Context, batch model.SampleBatch) error {
	if err := validation.Batch(batch); err != nil {
		return err
	}
	if _, err := s.stations.Get(ctx, batch.StationID); err != nil {
		return fmt.Errorf("sample station: %w", err)
	}
	if err := s.samples.BatchInsert(ctx, batch); err != nil {
		return fmt.Errorf("persist samples: %w", err)
	}
	s.metrics.Add("samples.accepted", int64(len(batch.Samples)))
	return nil
}

func (s *ObservatoryService) DrainResults(ctx context.Context, expected int) error {
	s.cacheMu.RLock()
	processor := s.processor
	s.cacheMu.RUnlock()
	if processor == nil {
		return nil
	}
	for i := 0; i < expected; i++ {
		select {
		case err := <-processor.Results():
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func (s *ObservatoryService) ProcessBatches(ctx context.Context, batches []model.SampleBatch) error {
	gate := make(chan struct{}, 1)
	for _, batch := range batches {
		gate <- struct{}{}
		err := ingest.ProcessSynchronously(ctx, s, []model.SampleBatch{batch})
		<-gate
		if err != nil {
			return err
		}
	}
	return nil
}
