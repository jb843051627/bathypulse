package handler

import "net/http"

func (s *Server) listAlerts(w http.ResponseWriter, r *http.Request) {
	items, err := s.service.ListAlerts(r.Context(), r.URL.Query().Get("station_id"), "")
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) ackAlert(w http.ResponseWriter, r *http.Request) {
	if err := s.service.AcknowledgeAlert(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"state": "acked"})
}
