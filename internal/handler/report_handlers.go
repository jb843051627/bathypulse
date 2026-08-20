package handler

import (
	"net/http"
	"time"
)

func (s *Server) summary(w http.ResponseWriter, r *http.Request) {
	value, err := s.service.Summary(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, value)
}

func (s *Server) timeline(w http.ResponseWriter, r *http.Request) {
	day, err := time.Parse("2006-01-02", r.URL.Query().Get("day"))
	if err != nil {
		writeError(w, err)
		return
	}
	items, err := s.service.Timeline(r.Context(), r.URL.Query().Get("station_id"), day)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}
