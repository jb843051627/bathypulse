package codec

import "time"

func parseTimestamp(value string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, value)
}

func FormatTimestamp(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}
