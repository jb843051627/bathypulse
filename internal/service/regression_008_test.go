package service

import (
	"context"
	"fmt"
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

func TestBug008_EventTimelineDoesNotExposeCacheSlice(t *testing.T) {
	svc := newRegressionService(t)
	addRegressionStation(t, svc, "station-008")
	for i := 0; i < 2; i++ {
		at := time.Date(2026, 8, 21, 7, i, 0, 0, time.UTC)
		if _, err := svc.OpenEvent(context.Background(), model.SeismicEvent{ID: fmt.Sprintf("event-008-%d", i), StationID: "station-008", Kind: "tremor", StartedAt: at, EndedAt: at.Add(time.Second), Magnitude: 2, Confidence: 0.9, State: model.EventOpen, Revision: 1}); err != nil {
			t.Fatal(err)
		}
	}
	first, err := svc.ListEvents(context.Background(), "station-008", "")
	if err != nil {
		t.Fatal(err)
	}
	first[0].Kind = "caller-value"
	second, err := svc.ListEvents(context.Background(), "station-008", "")
	if err != nil {
		t.Fatal(err)
	}
	if second[0].Kind == "caller-value" {
		t.Fatal("caller mutation changed event cache")
	}
}
