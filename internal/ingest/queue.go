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
	done    chan struct{}
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
	q := &Queue{jobs: make(chan Job, size), results: make(chan error, size), handler: handler, done: make(chan struct{})}
	q.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go q.worker()
	}
	return q
}

func (q *Queue) Submit(ctx context.Context, batch model.SampleBatch) error {
	job := Job{Context: ctx, Batch: batch}
	select {
	case q.jobs <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *Queue) worker() {
	defer q.wg.Done()
	for job := range q.jobs {
		result := q.handler.Handle(job.Context, job.Batch)
		// 结果通道可能因消费者取消或缓冲耗尽而无人接收。
		// 关闭流程触发后不应让被取消的 worker 因结果无人接收而拖住
		// 整个采集服务的关闭，故在 done 信号下放弃该结果。
		select {
		case q.results <- result:
		case <-q.done:
		}
	}
}

func (q *Queue) Results() <-chan error {
	return q.results
}

func (q *Queue) Close() {
	q.close.Do(func() {
		close(q.jobs)
		// 通知 worker：关闭流程已启动，无法投递的结果应被丢弃，
		// 避免 q.wg.Wait 因 worker 卡在 results 发送而永久挂起。
		close(q.done)
		q.wg.Wait()
		close(q.results)
	})
}
