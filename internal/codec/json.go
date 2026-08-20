package codec

import (
	"encoding/json"
	"fmt"

	"github.com/jb843051627/bathypulse/internal/model"
)

type Envelope struct {
	StationID  string `json:"station_id"`
	CapturedAt string `json:"captured_at"`
	Sequence   int64  `json:"sequence"`
	Rate       int    `json:"rate"`
	Payload    []byte `json:"payload"`
	Checksum   string `json:"checksum"`
}

func DecodeEnvelope(data []byte) (model.Waveform, error) {
	var envelope Envelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return model.Waveform{}, fmt.Errorf("decode envelope: %w", err)
	}
	at, err := parseTimestamp(envelope.CapturedAt)
	if err != nil {
		return model.Waveform{}, err
	}
	return model.Waveform{
		ID:         fmt.Sprintf("wave-%d", envelope.Sequence),
		StationID:  envelope.StationID,
		CapturedAt: at,
		DurationMs: len(envelope.Payload) * 10,
		SampleRate: envelope.Rate,
		Payload:    envelope.Payload,
		Checksum:   envelope.Checksum,
		State:      model.WaveformReceived,
		Sequence:   envelope.Sequence,
	}, nil
}

func EncodeEvent(event model.SeismicEvent) ([]byte, error) {
	return json.Marshal(event)
}
