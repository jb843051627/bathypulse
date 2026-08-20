package report

import (
	"sort"
	"time"

	"github.com/jb843051627/bathypulse/internal/codec"
)

type ReplayReport struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Records  int       `json:"records"`
	Stations []string  `json:"stations"`
	Kinds    []string  `json:"kinds"`
}

func BuildReplayReport(records []codec.ReplayRecord, start, end time.Time) ReplayReport {
	stations := make(map[string]struct{})
	kinds := make(map[string]struct{})
	for _, record := range records {
		if record.At.Before(start) || !record.At.Before(end) {
			continue
		}
		stations[record.StationID] = struct{}{}
		kinds[record.Kind] = struct{}{}
	}
	stationList := make([]string, 0, len(stations))
	for station := range stations {
		stationList = append(stationList, station)
	}
	kindList := make([]string, 0, len(kinds))
	for kind := range kinds {
		kindList = append(kindList, kind)
	}
	sort.Strings(stationList)
	sort.Strings(kindList)
	return ReplayReport{Start: start, End: end, Records: len(records), Stations: stationList, Kinds: kindList}
}

func ReplayCoverage(report ReplayReport) float64 {
	if report.Records == 0 {
		return 0
	}
	return float64(len(report.Stations)+len(report.Kinds)) / float64(report.Records)
}
