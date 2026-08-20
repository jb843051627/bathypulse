package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/validation"
)

func (s *ObservatoryService) QueueDispatch(ctx context.Context, job model.DispatchJob) (model.DispatchJob, error) {
	if err := validation.DispatchJob(job); err != nil {
		return job, err
	}
	if _, err := s.stations.Get(ctx, job.StationID); err != nil {
		return job, fmt.Errorf("dispatch station: %w", err)
	}
	if job.State == "" {
		job.State = model.DispatchQueued
	}
	if job.CreatedAt.IsZero() {
		job.CreatedAt = s.clock.Now()
	}
	job.UpdatedAt = job.CreatedAt
	if err := s.db.Dispatches().Create(ctx, job); err != nil {
		return job, fmt.Errorf("queue dispatch: %w", err)
	}
	return job, nil
}

func (s *ObservatoryService) StartDispatch(ctx context.Context, id, owner string) error {
	job, err := s.db.Dispatches().Get(ctx, id)
	if err != nil {
		return err
	}
	if !job.CanMoveTo(model.DispatchRunning) {
		return model.ErrConflict
	}
	leaseName := "dispatch:" + id
	acquired, err := s.db.Leases().Acquire(ctx, leaseName, owner, s.clock.Now().Add(2*time.Minute))
	if err != nil {
		return err
	}
	if !acquired {
		return fmt.Errorf("dispatch lease: %w", model.ErrConflict)
	}
	if err := s.db.Dispatches().Move(ctx, id, job.State, model.DispatchRunning, s.clock.Now(), job.Attempt); err != nil {
		_ = s.db.Leases().Release(ctx, leaseName, owner)
		return err
	}
	return nil
}

func (s *ObservatoryService) FinishDispatch(ctx context.Context, id, owner string, failed bool) error {
	job, err := s.db.Dispatches().Get(ctx, id)
	if err != nil {
		return err
	}
	next := model.DispatchCompleted
	if failed {
		next = model.DispatchFailed
	}
	if !job.CanMoveTo(next) {
		return model.ErrConflict
	}
	if err := s.db.Dispatches().Move(ctx, id, job.State, next, s.clock.Now(), job.Attempt+boolInt(failed)); err != nil {
		return err
	}
	return s.db.Leases().Release(ctx, "dispatch:"+id, owner)
}

func (s *ObservatoryService) RetryDispatch(ctx context.Context, id string) error {
	job, err := s.db.Dispatches().Get(ctx, id)
	if err != nil {
		return err
	}
	if !job.Retryable() {
		return model.ErrConflict
	}
	return s.db.Dispatches().Move(ctx, id, job.State, model.DispatchQueued, s.clock.Now(), job.Attempt+1)
}

func (s *ObservatoryService) QueuedDispatches(ctx context.Context, limit int) ([]model.DispatchJob, error) {
	return s.db.Dispatches().Queued(ctx, limit)
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
