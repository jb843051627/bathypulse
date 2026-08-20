package model

func AllowedStationTransition(from, to StationStatus) bool {
	if from == to {
		return true
	}
	if from == StationActive {
		return to == StationQuiet || to == StationOffline || to == StationRepair
	}
	if from == StationQuiet {
		return to == StationActive || to == StationOffline || to == StationRepair
	}
	if from == StationOffline {
		return to == StationActive || to == StationRepair
	}
	return from == StationRepair && to == StationActive
}

func AllowedAlertTransition(from, to AlertState) bool {
	if from == AlertPending {
		return to == AlertAcked || to == AlertClosed
	}
	if from == AlertAcked {
		return to == AlertClosed
	}
	return false
}
