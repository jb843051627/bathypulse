package model

import "time"

type DispatchBatch struct {
	ID        string        `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	Jobs      []DispatchJob `json:"jobs"`
	State     DispatchState `json:"state"`
}

func (b DispatchBatch) Validate() error {
	if b.ID == "" || len(b.Jobs) == 0 {
		return ErrInvalid
	}
	for _, job := range b.Jobs {
		if job.StationID == "" || job.Kind == "" {
			return ErrInvalid
		}
	}
	return nil
}

func (b DispatchBatch) Pending() int {
	count := 0
	for _, job := range b.Jobs {
		if job.State == DispatchQueued || job.State == DispatchRunning {
			count++
		}
	}
	return count
}

func (b DispatchBatch) FinishedAt() time.Time {
	finished := b.CreatedAt
	for _, job := range b.Jobs {
		if job.UpdatedAt.After(finished) {
			finished = job.UpdatedAt
		}
	}
	return finished
}
