package service

import (
	"context"
	"sync"
	"time"

	"github.com/jb843051627/bathypulse/internal/clock"
	"github.com/jb843051627/bathypulse/internal/ingest"
	"github.com/jb843051627/bathypulse/internal/metrics"
	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/store"
)

type ObservatoryService struct {
	db          *store.DB
	stations    store.StationRepository
	waveforms   store.WaveformRepository
	events      store.EventRepository
	alerts      store.AlertRepository
	maintenance store.MaintenanceRepository
	samples     store.SampleRepository
	queries     store.QueryRepository
	clock       clock.Clock
	metrics     *metrics.Registry
	cacheMu     sync.RWMutex
	waveCache   map[string][]model.Waveform
	eventCache  map[string][]model.SeismicEvent
	processor   *ingest.BatchProcessor
	policy      *policyState
}

func New(db *store.DB, now clock.Clock) *ObservatoryService {
	if now == nil {
		now = clock.System{}
	}
	return &ObservatoryService{
		db: db, stations: db.Stations(), waveforms: db.Waveforms(), events: db.Events(), alerts: db.Alerts(),
		maintenance: db.Maintenance(), samples: db.Samples(), queries: db.Queries(), clock: now,
		metrics: metrics.New(), waveCache: make(map[string][]model.Waveform), eventCache: make(map[string][]model.SeismicEvent),
		policy: newPolicyState(),
	}
}

func (s *ObservatoryService) Start() {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	if s.processor == nil {
		s.processor = ingest.NewProcessor(s)
	}
}

func (s *ObservatoryService) Close() {
	s.cacheMu.Lock()
	processor := s.processor
	s.processor = nil
	s.cacheMu.Unlock()
	if processor != nil {
		processor.Close()
	}
}

func (s *ObservatoryService) Handle(ctx context.Context, batch model.SampleBatch) error {
	return s.persistSamples(ctx, batch)
}

func (s *ObservatoryService) Health(ctx context.Context) error {
	return s.db.Ping(ctx)
}

func (s *ObservatoryService) QueueDepth() int {
	s.cacheMu.RLock()
	defer s.cacheMu.RUnlock()
	if s.processor == nil {
		return 0
	}
	return s.processor.Count()
}

func SeedStation(id, basin string) model.Station {
	return model.Station{ID: id, Code: id, Basin: basin, Latitude: -42.5, Longitude: 148.2, DepthMeters: 2800, Status: model.StationActive, LastSeen: time.Now().UTC(), Version: 1}
}
