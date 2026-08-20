package ingest

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type HealthMessage struct {
	StationID      string
	BatteryPercent float64
	SignalQuality  int
	ClockSkew      int64
	SampleCount    int
	At             time.Time
}

func ValidateHealth(ctx context.Context, message HealthMessage) (model.StationHealth, error) {
	if err := ctx.Err(); err != nil {
		return model.StationHealth{}, err
	}
	if message.StationID == "" || message.BatteryPercent < 0 || message.BatteryPercent > 100 || message.SignalQuality < 0 || message.SignalQuality > 100 {
		return model.StationHealth{}, fmt.Errorf("health message: %w", model.ErrInvalid)
	}
	return model.StationHealth{StationID: message.StationID, ObservedAt: message.At, BatteryPercent: message.BatteryPercent, SignalQuality: message.SignalQuality, ClockSkewMillis: message.ClockSkew, SampleCount: message.SampleCount}, nil
}

func HealthState(health model.StationHealth) model.StationStatus {
	if !health.Ready() {
		return model.StationOffline
	}
	if health.Score() < 0.5 {
		return model.StationQuiet
	}
	return model.StationActive
}
