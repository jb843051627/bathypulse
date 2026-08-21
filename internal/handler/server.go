package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jb843051627/bathypulse/internal/service"
)

type Server struct {
	service *service.ObservatoryService
}

func New(svc *service.ObservatoryService) *Server {
	return &Server{service: svc}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.dashboard)
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /v1/summary", s.summary)
	mux.HandleFunc("GET /v1/stations", s.listStations)
	mux.HandleFunc("POST /v1/stations", s.createStation)
	mux.HandleFunc("POST /v1/stations/{id}/heartbeat", s.heartbeat)
	mux.HandleFunc("POST /v1/stations/{id}/suspend", s.suspend)
	mux.HandleFunc("GET /v1/waveforms", s.listWaveforms)
	mux.HandleFunc("POST /v1/waveforms", s.createWaveform)
	mux.HandleFunc("GET /v1/events", s.listEvents)
	mux.HandleFunc("POST /v1/events/{id}/resolve", s.resolveEvent)
	mux.HandleFunc("GET /v1/alerts", s.listAlerts)
	mux.HandleFunc("POST /v1/alerts/{id}/ack", s.ackAlert)
	mux.HandleFunc("POST /v1/maintenance", s.createMaintenance)
	mux.HandleFunc("POST /v1/maintenance/{id}/complete", s.completeMaintenance)
	return mux
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}
