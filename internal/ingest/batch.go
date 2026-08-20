package ingest

import (
	"context"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

func ProcessSynchronously(ctx context.Context, handler Handler, batches []model.SampleBatch) error {
	for _, batch := range batches {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := handler.Handle(ctx, batch); err != nil {
			return err
		}
	}
	return nil
}

func NewBatch(stationID string, values []float64, at time.Time) model.SampleBatch {
	batch := model.SampleBatch{StationID: stationID, OpenedAt: at, Samples: make([]model.Sample, 0, len(values))}
	for i, value := range values {
		batch.Samples = append(batch.Samples, model.Sample{StationID: stationID, CapturedAt: at.Add(time.Duration(i) * time.Second), Sequence: int64(i + 1), Value: value, Quality: 90})
	}
	return batch
}
