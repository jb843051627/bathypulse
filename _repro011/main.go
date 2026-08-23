// 临时复现脚本（同 module 子包），验证 queue 并发关闭修复。跑完保留。
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jb843051627/bathypulse/internal/ingest"
	"github.com/jb843051627/bathypulse/internal/model"
)

// countingHandler 计数并快速返回，保证 worker 不永久阻塞，
// 这样 Close 的 wg.Wait 能正常结束。
type countingHandler struct {
	mu sync.Mutex
	n  int
}

func (h *countingHandler) Handle(ctx context.Context, b model.SampleBatch) error {
	h.mu.Lock()
	h.n++
	h.mu.Unlock()
	return nil
}

func main() {
	h := &countingHandler{}
	// 缓冲 1、worker 1：第 2 个 Submit 即会阻塞在 q.jobs <- job，与 Close 并发。
	q := ingest.NewQueue(1, 1, h)

	// 第 1 个 job 立即入队缓冲（worker 很可能还没取走）。
	_ = q.Submit(context.Background(), model.SampleBatch{
		StationID: "station-011",
		OpenedAt:  time.Now(),
		Samples:   []model.Sample{{StationID: "station-011", Sequence: 1, Value: 1}},
	})

	var (
		submitErr error
		submitWg  sync.WaitGroup
	)
	// 第 2 个 Submit 大概率阻塞在发送上（缓冲满 + worker 可能正忙），
	// 此时 Close 并发执行——修复前会 panic: send on closed channel。
	submitWg.Add(1)
	go func() {
		defer submitWg.Done()
		submitErr = q.Submit(context.Background(), model.SampleBatch{
			StationID: "station-011",
			OpenedAt:  time.Now(),
			Samples:   []model.Sample{{StationID: "station-011", Sequence: 99, Value: 1}},
		})
	}()

	q.Close()
	submitWg.Wait()

	afterErr := q.Submit(context.Background(), model.SampleBatch{
		StationID: "station-011",
		OpenedAt:  time.Now(),
		Samples:   []model.Sample{{StationID: "station-011", Sequence: 100, Value: 1}},
	})

	fmt.Printf("concurrent-submit-err = %v (is ErrQueueClosed=%v)\n", submitErr, errors.Is(submitErr, model.ErrQueueClosed))
	fmt.Printf("after-close-submit-err = %v (is ErrQueueClosed=%v)\n", afterErr, errors.Is(afterErr, model.ErrQueueClosed))
	fmt.Printf("handled = %d\n", h.n)

	if errors.Is(submitErr, model.ErrQueueClosed) && errors.Is(afterErr, model.ErrQueueClosed) {
		fmt.Println("RESULT: OK — no panic, ErrQueueClosed returned")
		return
	}
	fmt.Println("RESULT: FAIL")
}
