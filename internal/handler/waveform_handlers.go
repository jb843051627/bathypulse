package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jb843051627/bathypulse/internal/model"
)

func (s *Server) listWaveforms(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	items, err := s.service.ListWaveforms(r.Context(), r.URL.Query().Get("station_id"), limit)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) createWaveform(w http.ResponseWriter, r *http.Request) {
	var waveform model.Waveform
	if err := json.NewDecoder(r.Body).Decode(&waveform); err != nil {
		writeError(w, err)
		return
	}
	created, err := s.service.IngestWaveform(r.Context(), waveform)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}
