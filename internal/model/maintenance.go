package model

import "time"

func (m Maintenance) Contains(at time.Time) bool {
	return !at.Before(m.WindowStart) && at.Before(m.WindowEnd)
}

func (m Maintenance) CanMoveTo(next MaintenanceState) bool {
	if m.State == MaintenancePlanned {
		return next == MaintenanceRunning || next == MaintenanceCanceled
	}
	if m.State == MaintenanceRunning {
		return next == MaintenanceFinished || next == MaintenanceCanceled
	}
	return false
}
