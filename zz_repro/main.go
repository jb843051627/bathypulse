package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/store"
)

func main() {
	ctx := context.Background()
	dbPath := filepath.Join(os.TempDir(), fmt.Sprintf("bathypulse_repro_%d.db", time.Now().UnixNano()))
	db, err := store.Open(dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ev := model.SeismicEvent{
		ID: "ev-1", StationID: "st-1", Kind: "microseism",
		StartedAt: time.Now().UTC(), EndedAt: time.Now().UTC().Add(time.Minute),
		Magnitude: 3.2, Confidence: 0.9, State: model.EventOpen, Revision: 1,
	}
	if err := db.Events().Create(ctx, ev); err != nil {
		panic(err)
	}

	// 两个收敛任务都持有 revision=1 的快照，并发推进 open->review。
	errA := db.Events().UpdateState(ctx, "ev-1", model.EventOpen, model.EventReview, 1)
	errB := db.Events().UpdateState(ctx, "ev-1", model.EventOpen, model.EventReview, 1) // 陈旧 revision

	fmt.Println("task A (fresh revision 1):", errA)
	fmt.Println("task B (stale revision 1):", errB)
	fmt.Println("task B is ErrConflict:", errors.Is(errB, model.ErrConflict))

	got, _ := db.Events().Get(ctx, "ev-1")
	fmt.Printf("final: state=%s revision=%d\n", got.State, got.Revision)
}
