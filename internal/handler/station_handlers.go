package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/service"
)

func (s *Server) listStations(w http.ResponseWriter, r *http.Request) {
	items, err := s.service.ListStations(r.Context(), model.StationStatus(r.URL.Query().Get("status")))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) createStation(w http.ResponseWriter, r *http.Request) {
	var station model.Station
	if err := json.NewDecoder(r.Body).Decode(&station); err != nil {
		writeError(w, err)
		return
	}
	created, err := s.service.RegisterStation(r.Context(), station)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) heartbeat(w http.ResponseWriter, r *http.Request) {
	latency, _ := time.ParseDuration(r.URL.Query().Get("latency"))
	if err := s.service.UpdateHeartbeat(r.Context(), r.PathValue("id"), latency.Seconds()); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusNoContent, map[string]string{"status": "accepted"})
}

func (s *Server) suspend(w http.ResponseWriter, r *http.Request) {
	if err := s.service.SuspendStation(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "repair"})
}

var _ = service.SeedStation
