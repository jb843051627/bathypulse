package ingest

import (
	"context"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type WindowAssembler struct {
	stationID string
	maxSize   int
	values    []float64
	openedAt  time.Time
}

func NewWindowAssembler(stationID string, maxSize int) *WindowAssembler {
	if maxSize < 1 {
		maxSize = 64
	}
	return &WindowAssembler{stationID: stationID, maxSize: maxSize, values: make([]float64, 0, maxSize)}
}

func (a *WindowAssembler) Add(value float64, at time.Time) (model.SampleBatch, bool) {
	if a.openedAt.IsZero() {
		a.openedAt = at
	}
	a.values = append(a.values, value)
	if len(a.values) < a.maxSize {
		return model.SampleBatch{}, false
	}
	batch := NewBatch(a.stationID, a.values, a.openedAt)
	a.values = a.values[:0]
	a.openedAt = time.Time{}
	return batch, true
}

func (a *WindowAssembler) Flush() (model.SampleBatch, bool) {
	if len(a.values) == 0 {
		return model.SampleBatch{}, false
	}
	batch := NewBatch(a.stationID, a.values, a.openedAt)
	a.values = a.values[:0]
	a.openedAt = time.Time{}
	return batch, true
}

func RunAssembler(ctx context.Context, assembler *WindowAssembler, input <-chan float64, output chan<- model.SampleBatch) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case value, ok := <-input:
			if !ok {
				if batch, ready := assembler.Flush(); ready {
					output <- batch
				}
				return nil
			}
			if batch, ready := assembler.Add(value, time.Now().UTC()); ready {
				select {
				case output <- batch:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}
}
