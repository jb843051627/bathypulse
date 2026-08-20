package service

import (
	"context"
	"sync"

	"github.com/jb843051627/bathypulse/internal/model"
)

type policyState struct {
	once      sync.Once
	mu        sync.RWMutex
	threshold float64
	err       error
}

func newPolicyState() *policyState {
	return &policyState{threshold: 3.5}
}

func (s *ObservatoryService) EventRatioThreshold(ctx context.Context) (float64, error) {
	s.policy.once.Do(func() {
		if err := ctx.Err(); err != nil {
			s.policy.err = err
		}
	})
	s.policy.mu.RLock()
	defer s.policy.mu.RUnlock()
	return s.policy.threshold, s.policy.err
}

func (s *ObservatoryService) SetEventRatioThreshold(ctx context.Context, value float64) error {
	if value <= 0 {
		return model.ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	s.policy.mu.Lock()
	s.policy.threshold = value
	s.policy.mu.Unlock()
	return nil
}

func (s *ObservatoryService) PolicySnapshot() float64 {
	s.policy.mu.RLock()
	defer s.policy.mu.RUnlock()
	return s.policy.threshold
}
