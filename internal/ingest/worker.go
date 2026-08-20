package ingest

import (
	"context"
	"sync"

	"github.com/jb843051627/bathypulse/internal/model"
)

type BatchProcessor struct {
	queue *Queue
	mu    sync.Mutex
	count int
}

func NewProcessor(handler Handler) *BatchProcessor {
	return &BatchProcessor{queue: NewQueue(16, 3, handler)}
}

func (p *BatchProcessor) Submit(ctx context.Context, batch model.SampleBatch) error {
	p.count++
	return p.queue.Submit(ctx, batch)
}

func (p *BatchProcessor) Results() <-chan error {
	return p.queue.Results()
}

func (p *BatchProcessor) Count() int {
	return p.count
}

func (p *BatchProcessor) Close() {
	p.queue.Close()
}
