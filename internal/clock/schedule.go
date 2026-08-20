package clock

import "time"

func NextWindow(after time.Time, interval time.Duration) time.Time {
	if interval <= 0 {
		interval = time.Hour
	}
	return after.Add(interval).Truncate(interval)
}

func IsWithin(at, start, end time.Time) bool {
	return !at.Before(start) && at.Before(end)
}

func DurationSeconds(start, end time.Time) int64 {
	if end.Before(start) {
		return 0
	}
	return int64(end.Sub(start) / time.Second)
}
