package service

import (
	"context"
	"errors"
	"github.com/jb843051627/bathypulse/internal/clock"
	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/store"
	"github.com/jb843051627/bathypulse/internal/validation"
	"path/filepath"
	"testing"
	"time"
)

func newRegressionService(t *testing.T) *ObservatoryService {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	svc := New(db, clock.Fixed{Value: time.Date(2026, 8, 21, 8, 0, 0, 0, time.UTC)})
	t.Cleanup(func() {
		svc.Close()
		_ = db.Close()
	})
	return svc
}

func addRegressionStation(t *testing.T, svc *ObservatoryService, id string) {
	t.Helper()
	station := model.Station{ID: id, Code: id, Basin: "mid-ocean", Latitude: -10, Longitude: 140, DepthMeters: 2600, Status: model.StationActive, LastSeen: time.Date(2026, 8, 21, 7, 59, 0, 0, time.UTC), Version: 1}
	if _, err := svc.RegisterStation(context.Background(), station); err != nil {
		t.Fatal(err)
	}
}

func regressionWave(id, station string, at time.Time) model.Waveform {
	payload := []byte(id + "-payload")
	return model.Waveform{ID: id, StationID: station, CapturedAt: at, DurationMs: 1000, SampleRate: 100, Peak: 8, RMS: 1, Payload: payload, Checksum: validation.Checksum(payload), State: model.WaveformReceived, Sequence: 1}
}

func TestBug021_MaintenanceValidationPreservesSentinel(t *testing.T) {
	svc := newRegressionService(t)
	item := model.Maintenance{ID: "maint-021", StationID: "station-021", WindowStart: time.Date(2026, 8, 21, 11, 0, 0, 0, time.UTC), WindowEnd: time.Date(2026, 8, 21, 11, 0, 0, 0, time.UTC), Reason: "calibration"}
	_, err := svc.PlanMaintenance(context.Background(), item)
	if !errors.Is(err, model.ErrInvalid) {
		t.Fatalf("validation chain was lost: %v", err)
	}
}
