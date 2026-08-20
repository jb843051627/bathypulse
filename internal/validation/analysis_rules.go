package validation

import (
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

func AnalysisRun(run model.AnalysisRun) error {
	if run.ID == "" || run.StationID == "" || run.StartedAt.IsZero() {
		return fmt.Errorf("analysis identity: %w", model.ErrInvalid)
	}
	if run.Samples < 0 || run.Events < 0 || run.Score < 0 {
		return fmt.Errorf("analysis counters: %w", model.ErrInvalid)
	}
	return nil
}

func AnalysisWindow(start, end time.Time) error {
	if start.IsZero() || end.IsZero() || !end.After(start) {
		return fmt.Errorf("analysis window: %w", model.ErrInvalid)
	}
	if end.Sub(start) > 6*time.Hour {
		return fmt.Errorf("analysis window exceeds limit: %w", model.ErrInvalid)
	}
	return nil
}

func DispatchJob(job model.DispatchJob) error {
	if job.ID == "" || job.StationID == "" || job.Kind == "" {
		return fmt.Errorf("dispatch identity: %w", model.ErrInvalid)
	}
	if job.Attempt < 0 || job.Attempt > 3 {
		return fmt.Errorf("dispatch attempt: %w", model.ErrInvalid)
	}
	return nil
}

func Incident(item model.Incident) error {
	if item.ID == "" || item.StationID == "" || item.AlertID == "" || item.Summary == "" {
		return fmt.Errorf("incident fields: %w", model.ErrInvalid)
	}
	return nil
}
