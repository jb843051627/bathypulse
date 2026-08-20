package codec

import (
	"encoding/csv"
	"io"
	"strconv"

	"github.com/jb843051627/bathypulse/internal/model"
)

func WriteTimeline(w io.Writer, items []model.TimelineItem) error {
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"at", "kind", "entity_id", "station", "summary"}); err != nil {
		return err
	}
	for _, item := range items {
		if err := writer.Write([]string{
			item.At.Format("2006-01-02T15:04:05Z07:00"),
			item.Kind,
			item.EntityID,
			item.Station,
			item.Summary,
		}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func WriteSample(w io.Writer, sample model.Sample) error {
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{sample.StationID, sample.CapturedAt.Format(timeLayout), strconv.FormatInt(sample.Sequence, 10), strconv.FormatFloat(sample.Value, 'f', 4, 64), strconv.Itoa(sample.Quality)}); err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}

const timeLayout = "2006-01-02T15:04:05Z07:00"
