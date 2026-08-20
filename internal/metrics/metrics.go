package metrics

import "sync"

type Registry struct {
	mu       sync.RWMutex
	counters map[string]int64
	gauges   map[string]float64
}

func New() *Registry {
	return &Registry{counters: make(map[string]int64), gauges: make(map[string]float64)}
}

func (r *Registry) Add(name string, delta int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counters[name] += delta
}

func (r *Registry) Set(name string, value float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.gauges[name] = value
}

func (r *Registry) ObserveStation(stationID string, value float64) {
	r.Add("station."+stationID+".heartbeat", 1)
	r.Set("station."+stationID+".latency", value)
}

func (r *Registry) Snapshot() map[string]float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]float64, len(r.counters)+len(r.gauges))
	for name, value := range r.counters {
		out[name] = float64(value)
	}
	for name, value := range r.gauges {
		out[name] = value
	}
	return out
}

func (r *Registry) Count(name string) int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.counters[name]
}
