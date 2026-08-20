package ingest

import (
	"context"
	"sync"
	"time"
)

type Schedule struct {
	Name     string
	Interval time.Duration
	Action   func(context.Context) error
}

type Scheduler struct {
	mu      sync.Mutex
	entries map[string]Schedule
}

func NewScheduler() *Scheduler {
	return &Scheduler{entries: make(map[string]Schedule)}
}

func (s *Scheduler) Add(item Schedule) {
	s.mu.Lock()
	s.entries[item.Name] = item
	s.mu.Unlock()
}

func (s *Scheduler) Remove(name string) {
	s.mu.Lock()
	delete(s.entries, name)
	s.mu.Unlock()
}

func (s *Scheduler) Run(ctx context.Context, name string) error {
	s.mu.Lock()
	item, ok := s.entries[name]
	s.mu.Unlock()
	if !ok {
		return context.Canceled
	}
	ticker := time.NewTicker(item.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := item.Action(ctx); err != nil {
				return err
			}
		}
	}
}
