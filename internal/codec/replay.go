package codec

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type ReplayRecord struct {
	At        time.Time       `json:"at"`
	StationID string          `json:"station_id"`
	Kind      string          `json:"kind"`
	Payload   json.RawMessage `json:"payload"`
}

func DecodeReplay(data []byte) (ReplayRecord, error) {
	var record ReplayRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return record, fmt.Errorf("replay record: %w", err)
	}
	if record.At.IsZero() || record.StationID == "" || len(record.Payload) == 0 {
		return ReplayRecord{}, model.ErrInvalid
	}
	return record, nil
}

func EncodeReplay(record ReplayRecord) ([]byte, error) {
	return json.Marshal(record)
}

func ReplayWindow(records []ReplayRecord, start, end time.Time) []ReplayRecord {
	result := make([]ReplayRecord, 0, len(records))
	for _, record := range records {
		if !record.At.Before(start) && record.At.Before(end) {
			result = append(result, record)
		}
	}
	return result
}
