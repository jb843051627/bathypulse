package model

import "time"

type DispatchState string

const (
	DispatchQueued    DispatchState = "queued"
	DispatchRunning   DispatchState = "running"
	DispatchCompleted DispatchState = "completed"
	DispatchFailed    DispatchState = "failed"
)

type DispatchJob struct {
	ID        string        `json:"id"`
	StationID string        `json:"station_id"`
	Kind      string        `json:"kind"`
	State     DispatchState `json:"state"`
	Attempt   int           `json:"attempt"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Payload   string        `json:"payload"`
}

func (j DispatchJob) CanMoveTo(next DispatchState) bool {
	if j.State == DispatchQueued {
		return next == DispatchRunning || next == DispatchFailed
	}
	if j.State == DispatchRunning {
		return next == DispatchCompleted || next == DispatchFailed
	}
	return false
}

func (j DispatchJob) Retryable() bool {
	return j.State == DispatchFailed && j.Attempt < 3
}
