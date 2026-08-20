package model

func (a Alert) IsActive() bool {
	return a.State == AlertPending || a.State == AlertAcked
}

func (a Alert) NeedsEscalation() bool {
	return a.State == AlertPending && a.Level != AlertCritical
}
