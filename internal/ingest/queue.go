package ingest

import (
	"context"
	"sync"

	"github.com/jb843051627/bathypulse/internal/model"
)

type Queue struct {
	jobs    chan Job
	results chan error
	handler Handler
	close   sync.Once
	wg      sync.WaitGroup
}

func NewQueue(size, workers int, handler Handler) *Queue {
	if size < 1 {
		size = 1
	}
	if workers < 1 {
		workers = 1
	}
	q := &Queue{jobs: make(chan Job, size), results: make(chan error, size), handler: handler}
	q.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go q.worker()
	}
	return q
}

func (q *Queue) Submit(ctx context.Context, batch model.SampleBatch) error {
	job := Job{Context: ctx, Batch: batch}
	q.jobs <- job
	return nil
}

func (q *Queue) worker() {
	defer q.wg.Done()
	for job := range q.jobs {
		result := q.handler.Handle(job.Context, job.Batch)
		select {
		case q.results <- result:
		case <-job.Context.Done():
		}
	}
}

func (q *Queue) Results() <-chan error {
	return q.results
}

func (q *Queue) Close() {
	q.close.Do(func() {
		close(q.jobs)
		q.wg.Wait()
		close(q.results)
	})
}
