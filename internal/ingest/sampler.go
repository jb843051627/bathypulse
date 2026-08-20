package ingest

import (
	"context"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type Sampler struct {
	interval time.Duration
	clock    func() time.Time
}

func NewSampler(interval time.Duration, now func() time.Time) *Sampler {
	if interval <= 0 {
		interval = time.Second
	}
	return &Sampler{interval: interval, clock: now}
}

func (s *Sampler) Run(ctx context.Context, stationID string, emit func(model.Sample) error) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	var sequence int64
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case at := <-ticker.C:
			sequence++
			if err := emit(model.Sample{StationID: stationID, CapturedAt: at, Sequence: sequence, Value: 0, Quality: 100}); err != nil {
				return err
			}
		}
	}
}
