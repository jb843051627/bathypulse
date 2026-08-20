package model

import "time"

type WaveformSegment struct {
	ID         string    `json:"id"`
	WaveformID string    `json:"waveform_id"`
	OffsetMs   int       `json:"offset_ms"`
	DurationMs int       `json:"duration_ms"`
	Peak       float64   `json:"peak"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s WaveformSegment) EndOffset() int {
	return s.OffsetMs + s.DurationMs
}

func (s WaveformSegment) Overlaps(start, end int) bool {
	return s.OffsetMs < end && s.EndOffset() > start
}

func (s WaveformSegment) IsImpulse() bool {
	return s.DurationMs <= 500 && s.Peak >= 1
}
