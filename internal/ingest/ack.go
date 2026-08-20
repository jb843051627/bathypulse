package ingest

import (
	"context"
	"sync"
)

type AckTracker struct {
	mu      sync.Mutex
	pending map[int64]struct{}
}

func NewAckTracker() *AckTracker {
	return &AckTracker{pending: make(map[int64]struct{})}
}

func (t *AckTracker) Add(sequence int64) {
	t.mu.Lock()
	t.pending[sequence] = struct{}{}
	t.mu.Unlock()
}

func (t *AckTracker) Confirm(sequence int64) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.pending[sequence]; !ok {
		return false
	}
	delete(t.pending, sequence)
	return true
}

func (t *AckTracker) Wait(ctx context.Context, sequence int64) error {
	for {
		if !t.Confirm(sequence) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	}
}
