package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jb843051627/bathypulse/internal/model"
)

func (s *Server) createMaintenance(w http.ResponseWriter, r *http.Request) {
	var item model.Maintenance
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, err)
		return
	}
	created, err := s.service.PlanMaintenance(r.Context(), item)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) completeMaintenance(w http.ResponseWriter, r *http.Request) {
	if err := s.service.CompleteMaintenance(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"state": "finished"})
}
