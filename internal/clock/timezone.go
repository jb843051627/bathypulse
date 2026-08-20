package clock

import "time"

var basinZones = map[string]string{
	"north-trench": "Pacific/Auckland",
	"mid-ocean":    "UTC",
	"south-ridge":  "Australia/Hobart",
}

func LocationForBasin(basin string) *time.Location {
	name := basinZones[basin]
	if name == "" {
		return time.UTC
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC
	}
	return loc
}

func LocalDay(at time.Time, basin string) string {
	return at.In(LocationForBasin(basin)).Format("2006-01-02")
}
