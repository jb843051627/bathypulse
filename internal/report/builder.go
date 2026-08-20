package report

import (
	"sort"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type Builder struct{}

func (Builder) Timeline(events []model.SeismicEvent, alerts []model.Alert, basin string) []model.TimelineItem {
	items := make([]model.TimelineItem, 0, len(events)+len(alerts))
	for _, event := range events {
		items = append(items, model.TimelineItem{At: event.StartedAt, Kind: "event", EntityID: event.ID, Station: event.StationID, Summary: event.Kind})
	}
	for _, alert := range alerts {
		items = append(items, model.TimelineItem{At: alert.CreatedAt, Kind: "alert", EntityID: alert.ID, Station: alert.StationID, Summary: alert.Message})
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].At.Before(items[j].At) })
	return items
}

func (Builder) DayBounds(day time.Time, loc *time.Location) (time.Time, time.Time) {
	local := day.In(loc)
	start := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
	return start.UTC(), start.Add(24 * time.Hour).UTC()
}
