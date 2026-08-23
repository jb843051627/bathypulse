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
	mu      sync.RWMutex
	closed  bool
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

// Submit enqueues a batch for asynchronous handling. It is safe to call
// concurrently with Close: once the queue has been closed, Submit returns
// model.ErrQueueClosed instead of panicking on a send to a closed channel.
//
// The read lock is held across both the closed-flag check and the send so
// that Close cannot close the jobs channel in between the two (which would
// otherwise race and panic). Close acquires the write lock, so it is forced
// to wait until any in-flight Submit has finished sending before it closes
// the channel.
func (q *Queue) Submit(ctx context.Context, batch model.SampleBatch) error {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if q.closed {
		return model.ErrQueueClosed
	}
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
		q.mu.Lock()
		q.closed = true
		close(q.jobs)
		q.mu.Unlock()
		q.wg.Wait()
		close(q.results)
	})
}
