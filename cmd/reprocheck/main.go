package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jb843051627/bathypulse/internal/ingest"
	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/service"
	"github.com/jb843051627/bathypulse/internal/store"
)

func main() {
	dir := os.TempDir()
	dbPath := filepath.Join(dir, fmt.Sprintf("repro016-%d.db", time.Now().UnixNano()))
	_ = os.Remove(dbPath)
	db, err := store.Open(dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	station := model.Station{
		ID: "station-016", Code: "station-016", Basin: "bathy",
		Latitude: -42.5, Longitude: 148.2, DepthMeters: 2800,
		Status: model.StationActive, LastSeen: time.Now().UTC(), Version: 1,
	}
	if err := db.Stations().Create(context.Background(), station); err != nil {
		fmt.Println("create station err:", err)
		return
	}

	svc := service.New(db, nil)

	at := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)

	// 第一批：序列 1,2 成功落库
	b1 := ingest.NewBatch("station-016", []float64{1.0, 2.0}, at)
	if err := svc.SubmitSamples(context.Background(), b1); err != nil {
		fmt.Println("batch1 err:", err)
	} else {
		fmt.Println("batch1 ok")
	}

	// 场景A：冲突在【最后一条】。前两条(seq 3,4)在 tx 内成功，第三条(seq 1,与已有冲突)
	// 应触发整体回滚。若回滚失效，3,4 残留，Latest 会读到 4 条。
	b2 := model.SampleBatch{StationID: "station-016", OpenedAt: at.Add(time.Minute)}
	b2.Samples = []model.Sample{
		{StationID: "station-016", CapturedAt: at.Add(time.Minute), Sequence: 3, Value: 30.0, Quality: 90},
		{StationID: "station-016", CapturedAt: at.Add(time.Minute + time.Second), Sequence: 4, Value: 40.0, Quality: 90},
		{StationID: "station-016", CapturedAt: at.Add(time.Minute + 2*time.Second), Sequence: 1, Value: 999.0, Quality: 90}, // 冲突
	}
	if err := svc.SubmitSamples(context.Background(), b2); err != nil {
		fmt.Println("batch2 err (expected):", err)
	} else {
		fmt.Println("batch2 ok (unexpected!)")
	}

	latest, _ := db.Samples().Latest(context.Background(), "station-016", 10)
	fmt.Printf("[A] After failed batch2 (conflict on last row): Latest returned %d samples (expect 2 if rollback worked; 4 if leaked)\n", len(latest))
	for _, s := range latest {
		fmt.Printf("  seq=%d value=%v\n", s.Sequence, s.Value)
	}

	// 场景B：ctx 在 BatchInsert 循环【中途】取消。
	// 用大量样本(2000 条) + 很短 timeout(3ms)，确保取消发生在循环中途。
	// 修复前：execTx 用 context.Background() 不感知取消，全部插入成功并 Commit，返回 nil —— 调用方无错误但 2000 条全落库。
	// 修复后：execTx 用调用方 ctx，取消中断 exec，fn 返回 err 触发 Rollback，返回错误且无新增。
	const N = 2000
	toCtx, toCancel := context.WithTimeout(context.Background(), 3*time.Millisecond)
	defer toCancel()
	b3 := model.SampleBatch{StationID: "station-016", OpenedAt: at.Add(2 * time.Minute)}
	b3.Samples = make([]model.Sample, 0, N)
	for i := 0; i < N; i++ {
		b3.Samples = append(b3.Samples, model.Sample{
			StationID:  "station-016",
			CapturedAt: at.Add(2*time.Minute + time.Duration(i)*time.Millisecond),
			Sequence:   int64(1000 + i),
			Value:      float64(i),
			Quality:    90,
		})
	}
	errB := svc.SubmitSamples(toCtx, b3)
	if errB != nil {
		fmt.Printf("batch3 (mid-cancel ctx) err (expected after fix): %v\n", errB)
	} else {
		fmt.Println("batch3 ok (bug: cancel ignored, all inserted)")
	}

	latest2, _ := db.Samples().Latest(context.Background(), "station-016", 10)
	fmt.Printf("[B] After mid-cancel batch3: Latest returned %d samples (expect 2 if rollback worked; 5 if leaked)\n", len(latest2))
	for _, s := range latest2 {
		fmt.Printf("  seq=%d value=%v\n", s.Sequence, s.Value)
	}
}
