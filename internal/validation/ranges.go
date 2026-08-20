package validation

import (
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

func Window(start, end time.Time) error {
	if start.IsZero() || end.IsZero() || !end.After(start) {
		return fmt.Errorf("invalid time window: %w", model.ErrInvalid)
	}
	if end.Sub(start) > 30*time.Minute {
		return fmt.Errorf("time window too long: %w", model.ErrInvalid)
	}
	return nil
}

func EventThreshold(peak, rms float64) (float64, bool) {
	if peak <= 0 || rms <= 0 || peak < rms {
		return 0, false
	}
	value := peak / rms
	return value, value >= 3.5
}

func QualityWeight(quality int) float64 {
	if quality <= 0 {
		return 0
	}
	if quality >= 100 {
		return 1
	}
	return float64(quality) / 100
}
