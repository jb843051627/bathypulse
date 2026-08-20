package ingest

import (
	"context"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type Job struct {
	QueuedAt time.Time
	Context  context.Context
	Batch    model.SampleBatch
	Created  time.Time
}

type Handler interface {
	Handle(context.Context, model.SampleBatch) error
}
