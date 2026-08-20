package ingest

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type Decoder struct{}

func (Decoder) DecodeSamples(stationID string, values []float64, start time.Time) (model.SampleBatch, error) {
	if stationID == "" || len(values) == 0 || start.IsZero() {
		return model.SampleBatch{}, model.ErrInvalid
	}
	batch := model.SampleBatch{StationID: stationID, OpenedAt: start, Samples: make([]model.Sample, 0, len(values))}
	for i, value := range values {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return model.SampleBatch{}, fmt.Errorf("sample %d: %w", i, model.ErrInvalid)
		}
		batch.Samples = append(batch.Samples, model.Sample{StationID: stationID, CapturedAt: start.Add(time.Duration(i) * time.Second), Sequence: int64(i + 1), Value: value, Quality: 100})
	}
	return batch, nil
}

func EncodeFloat(value float64) []byte {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, math.Float64bits(value))
	return data
}
