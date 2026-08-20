package model

import (
	"errors"
	"time"
)

var (
	ErrNotFound    = errors.New("record not found")
	ErrConflict    = errors.New("state conflict")
	ErrInvalid     = errors.New("invalid observatory data")
	ErrChecksum    = errors.New("waveform checksum mismatch")
	ErrQueueClosed = errors.New("ingest queue closed")
)

type StationStatus string

const (
	StationActive  StationStatus = "active"
	StationQuiet   StationStatus = "quiet"
	StationOffline StationStatus = "offline"
	StationRepair  StationStatus = "repair"
)

type WaveformState string

const (
	WaveformReceived  WaveformState = "received"
	WaveformProcessed WaveformState = "processed"
	WaveformRejected  WaveformState = "rejected"
)

type EventState string

const (
	EventOpen     EventState = "open"
	EventReview   EventState = "review"
	EventResolved EventState = "resolved"
	EventArchived EventState = "archived"
)

type AlertLevel string

const (
	AlertInfo     AlertLevel = "info"
	AlertWatch    AlertLevel = "watch"
	AlertCritical AlertLevel = "critical"
)

type AlertState string

const (
	AlertPending AlertState = "pending"
	AlertAcked   AlertState = "acked"
	AlertClosed  AlertState = "closed"
)

type MaintenanceState string

const (
	MaintenancePlanned  MaintenanceState = "planned"
	MaintenanceRunning  MaintenanceState = "running"
	MaintenanceFinished MaintenanceState = "finished"
	MaintenanceCanceled MaintenanceState = "canceled"
)

type Station struct {
	ID          string        `json:"id"`
	Code        string        `json:"code"`
	Basin       string        `json:"basin"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
	DepthMeters int           `json:"depth_meters"`
	Status      StationStatus `json:"status"`
	LastSeen    time.Time     `json:"last_seen"`
	Version     int           `json:"version"`
}

type Waveform struct {
	ID         string        `json:"id"`
	StationID  string        `json:"station_id"`
	CapturedAt time.Time     `json:"captured_at"`
	DurationMs int           `json:"duration_ms"`
	SampleRate int           `json:"sample_rate"`
	Peak       float64       `json:"peak"`
	RMS        float64       `json:"rms"`
	Payload    []byte        `json:"payload"`
	Checksum   string        `json:"checksum"`
	State      WaveformState `json:"state"`
	Sequence   int64         `json:"sequence"`
}

type Sample struct {
	StationID  string    `json:"station_id"`
	CapturedAt time.Time `json:"captured_at"`
	Sequence   int64     `json:"sequence"`
	Value      float64   `json:"value"`
	Quality    int       `json:"quality"`
}

type SampleBatch struct {
	StationID string    `json:"station_id"`
	OpenedAt  time.Time `json:"opened_at"`
	Samples   []Sample  `json:"samples"`
}

type SeismicEvent struct {
	ID            string     `json:"id"`
	StationID     string     `json:"station_id"`
	Kind          string     `json:"kind"`
	StartedAt     time.Time  `json:"started_at"`
	EndedAt       time.Time  `json:"ended_at"`
	Magnitude     float64    `json:"magnitude"`
	Confidence    float64    `json:"confidence"`
	State         EventState `json:"state"`
	WaveformCount int        `json:"waveform_count"`
	Revision      int        `json:"revision"`
}

type Alert struct {
	ID        string     `json:"id"`
	EventID   string     `json:"event_id"`
	StationID string     `json:"station_id"`
	Level     AlertLevel `json:"level"`
	State     AlertState `json:"state"`
	Message   string     `json:"message"`
	CreatedAt time.Time  `json:"created_at"`
	AckedAt   *time.Time `json:"acked_at,omitempty"`
	Revision  int        `json:"revision"`
}

type Maintenance struct {
	ID          string           `json:"id"`
	StationID   string           `json:"station_id"`
	WindowStart time.Time        `json:"window_start"`
	WindowEnd   time.Time        `json:"window_end"`
	State       MaintenanceState `json:"state"`
	Reason      string           `json:"reason"`
	Revision    int              `json:"revision"`
}

type TimelineItem struct {
	At       time.Time `json:"at"`
	Kind     string    `json:"kind"`
	EntityID string    `json:"entity_id"`
	Station  string    `json:"station"`
	Summary  string    `json:"summary"`
}

type ObservatorySummary struct {
	GeneratedAt    time.Time `json:"generated_at"`
	StationCount   int       `json:"station_count"`
	ActiveStations int       `json:"active_stations"`
	OpenEvents     int       `json:"open_events"`
	PendingAlerts  int       `json:"pending_alerts"`
	Samples        int       `json:"samples"`
}

func CloneBytes(src []byte) []byte {
	if src == nil {
		return nil
	}
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

func CloneWaveform(src Waveform) Waveform {
	src.Payload = CloneBytes(src.Payload)
	return src
}

func CloneWaveforms(src []Waveform) []Waveform {
	dst := make([]Waveform, len(src))
	for i, item := range src {
		dst[i] = CloneWaveform(item)
	}
	return dst
}

func CloneEvents(src []SeismicEvent) []SeismicEvent {
	dst := make([]SeismicEvent, len(src))
	copy(dst, src)
	return dst
}
