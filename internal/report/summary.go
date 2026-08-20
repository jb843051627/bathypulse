package report

import (
	"sort"

	"github.com/jb843051627/bathypulse/internal/model"
)

func RankStations(stations []model.Station, samples map[string]int) []model.Station {
	result := append([]model.Station(nil), stations...)
	sort.SliceStable(result, func(i, j int) bool { return samples[result[i].ID] > samples[result[j].ID] })
	return result
}

func EventKinds(events []model.SeismicEvent) map[string]int {
	counts := make(map[string]int)
	for _, event := range events {
		counts[event.Kind]++
	}
	return counts
}
