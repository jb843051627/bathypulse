package codec

import (
	"encoding/json"
	"fmt"
	"time"
)

type DiagnosticNote struct {
	StationID string    `json:"station_id"`
	At        time.Time `json:"at"`
	Category  string    `json:"category"`
	Message   string    `json:"message"`
	Values    []float64 `json:"values"`
}

func EncodeDiagnostic(note DiagnosticNote) ([]byte, error) {
	if note.StationID == "" || note.Category == "" || note.At.IsZero() {
		return nil, fmt.Errorf("diagnostic fields are incomplete")
	}
	return json.Marshal(note)
}

func DecodeDiagnostic(data []byte) (DiagnosticNote, error) {
	var note DiagnosticNote
	if err := json.Unmarshal(data, &note); err != nil {
		return note, err
	}
	if note.StationID == "" || note.Category == "" {
		return DiagnosticNote{}, fmt.Errorf("diagnostic fields are incomplete")
	}
	return note, nil
}
