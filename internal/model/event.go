package model

import "time"

func (e SeismicEvent) Duration() time.Duration {
	if e.EndedAt.Before(e.StartedAt) {
		return 0
	}
	return e.EndedAt.Sub(e.StartedAt)
}

func (e SeismicEvent) IsTerminal() bool {
	return e.State == EventResolved || e.State == EventArchived
}

func (e SeismicEvent) CanMoveTo(next EventState) bool {
	if e.State == EventOpen {
		return next == EventReview || next == EventResolved
	}
	if e.State == EventReview {
		return next == EventOpen || next == EventResolved
	}
	if e.State == EventResolved {
		return next == EventArchived
	}
	return false
}
