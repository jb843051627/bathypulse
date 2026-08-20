package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

func (s *ObservatoryService) AcquireLease(ctx context.Context, name, owner string, ttl time.Duration) error {
	if name == "" || owner == "" || ttl <= 0 || ttl > time.Hour {
		return model.ErrInvalid
	}
	ok, err := s.db.Leases().Acquire(ctx, name, owner, s.clock.Now().Add(ttl))
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("lease busy: %w", model.ErrConflict)
	}
	return nil
}

func (s *ObservatoryService) ReleaseLease(ctx context.Context, name, owner string) error {
	if name == "" || owner == "" {
		return model.ErrInvalid
	}
	return s.db.Leases().Release(ctx, name, owner)
}

func (s *ObservatoryService) LeaseOwner(ctx context.Context, name string) (string, error) {
	lease, err := s.db.Leases().Get(ctx, name)
	if err != nil {
		return "", err
	}
	if !lease.ExpiresAt.After(s.clock.Now()) {
		return "", model.ErrNotFound
	}
	return lease.Owner, nil
}
