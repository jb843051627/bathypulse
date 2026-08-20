package handler

import (
	"net/http"

	"github.com/jb843051627/bathypulse/internal/model"
)

func (s *Server) listEvents(w http.ResponseWriter, r *http.Request) {
	items, err := s.service.ListEvents(r.Context(), r.URL.Query().Get("station_id"), model.EventState(r.URL.Query().Get("state")))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) resolveEvent(w http.ResponseWriter, r *http.Request) {
	if err := s.service.MoveEvent(r.Context(), r.PathValue("id"), model.EventResolved); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"state": string(model.EventResolved)})
}
