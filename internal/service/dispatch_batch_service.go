package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bathypulse/internal/model"
)

func (s *ObservatoryService) QueueDispatchBatch(ctx context.Context, batch model.DispatchBatch) error {
	if err := batch.Validate(); err != nil {
		return err
	}
	for _, job := range batch.Jobs {
		if _, err := s.QueueDispatch(ctx, job); err != nil {
			return fmt.Errorf("queue dispatch batch: %w", err)
		}
	}
	return nil
}

func (s *ObservatoryService) DispatchBatchPending(batch model.DispatchBatch) int {
	return batch.Pending()
}

func (s *ObservatoryService) FinishDispatchBatch(ctx context.Context, batch model.DispatchBatch, owner string, failed bool) error {
	if err := batch.Validate(); err != nil {
		return err
	}
	for _, job := range batch.Jobs {
		if job.State != model.DispatchRunning {
			continue
		}
		if err := s.FinishDispatch(ctx, job.ID, owner, failed); err != nil {
			return err
		}
	}
	return nil
}

func (s *ObservatoryService) RetryFailedDispatches(ctx context.Context, jobs []model.DispatchJob) (int, error) {
	count := 0
	for _, job := range jobs {
		if !job.Retryable() {
			continue
		}
		if err := s.RetryDispatch(ctx, job.ID); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}
