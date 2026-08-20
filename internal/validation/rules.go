package validation

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

func Station(s model.Station) error {
	if strings.TrimSpace(s.ID) == "" || strings.TrimSpace(s.Code) == "" {
		return fmt.Errorf("station identity: %w", model.ErrInvalid)
	}
	if s.Basin == "" || s.DepthMeters <= 0 {
		return fmt.Errorf("station geometry: %w", model.ErrInvalid)
	}
	if s.Latitude < -90 || s.Latitude > 90 || s.Longitude < -180 || s.Longitude > 180 {
		return fmt.Errorf("station coordinates: %w", model.ErrInvalid)
	}
	return nil
}

func Waveform(w model.Waveform) error {
	if w.StationID == "" || w.ID == "" || w.DurationMs <= 0 {
		return fmt.Errorf("waveform identity: %w", model.ErrInvalid)
	}
	if w.SampleRate < 10 || w.SampleRate > 10000 {
		return fmt.Errorf("waveform sample rate: %w", model.ErrInvalid)
	}
	if w.CapturedAt.IsZero() || len(w.Payload) == 0 {
		return fmt.Errorf("waveform payload: %w", model.ErrInvalid)
	}
	return nil
}

func Event(e model.SeismicEvent) error {
	if e.ID == "" || e.StationID == "" || e.Kind == "" {
		return fmt.Errorf("event identity: %w", model.ErrInvalid)
	}
	if e.StartedAt.IsZero() || e.EndedAt.Before(e.StartedAt) {
		return fmt.Errorf("event window: %w", model.ErrInvalid)
	}
	if e.Confidence < 0 || e.Confidence > 1 || e.Magnitude < 0 {
		return fmt.Errorf("event score: %w", model.ErrInvalid)
	}
	return nil
}

func Maintenance(m model.Maintenance) error {
	if m.ID == "" || m.StationID == "" || strings.TrimSpace(m.Reason) == "" {
		return fmt.Errorf("maintenance identity: %w", model.ErrInvalid)
	}
	if m.WindowStart.IsZero() || !m.WindowEnd.After(m.WindowStart) {
		return fmt.Errorf("maintenance window: %w", model.ErrInvalid)
	}
	return nil
}

func Batch(b model.SampleBatch) error {
	if b.StationID == "" || b.OpenedAt.IsZero() || len(b.Samples) == 0 {
		return fmt.Errorf("sample batch: %w", model.ErrInvalid)
	}
	previous := time.Time{}
	for _, sample := range b.Samples {
		if sample.StationID != b.StationID || sample.Quality < 0 || sample.Quality > 100 {
			return fmt.Errorf("sample member: %w", model.ErrInvalid)
		}
		if !previous.IsZero() && sample.CapturedAt.Before(previous) {
			return fmt.Errorf("sample order: %w", model.ErrInvalid)
		}
		previous = sample.CapturedAt
	}
	return nil
}

func Checksum(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func MatchesChecksum(payload []byte, expected string) bool {
	return Checksum(payload) == strings.ToLower(strings.TrimSpace(expected))
}
