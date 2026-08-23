package handler

import (
	"github.com/jb843051627/bathypulse/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBug028_NotFoundErrorsUseNotFoundStatus(t *testing.T) {
	recorder := httptest.NewRecorder()
	writeError(recorder, model.ErrNotFound)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d", recorder.Code)
	}
}
