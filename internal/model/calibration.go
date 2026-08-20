package model

import (
	"fmt"
	"sort"
)

type CalibrationPoint struct {
	Input  float64 `json:"input"`
	Output float64 `json:"output"`
}

type CalibrationProfile struct {
	ID        string             `json:"id"`
	StationID string             `json:"station_id"`
	Version   int                `json:"version"`
	Points    []CalibrationPoint `json:"points"`
	Enabled   bool               `json:"enabled"`
}

func (p CalibrationProfile) Validate() error {
	if p.ID == "" || p.StationID == "" || len(p.Points) < 2 {
		return fmt.Errorf("calibration profile: %w", ErrInvalid)
	}
	points := append([]CalibrationPoint(nil), p.Points...)
	sort.Slice(points, func(i, j int) bool { return points[i].Input < points[j].Input })
	for i := 1; i < len(points); i++ {
		if points[i].Input == points[i-1].Input {
			return fmt.Errorf("calibration input: %w", ErrInvalid)
		}
	}
	return nil
}

func (p CalibrationProfile) Apply(input float64) float64 {
	if len(p.Points) == 0 {
		return input
	}
	points := append([]CalibrationPoint(nil), p.Points...)
	sort.Slice(points, func(i, j int) bool { return points[i].Input < points[j].Input })
	if input <= points[0].Input {
		return points[0].Output
	}
	for i := 1; i < len(points); i++ {
		left, right := points[i-1], points[i]
		if input <= right.Input {
			ratio := (input - left.Input) / (right.Input - left.Input)
			return left.Output + ratio*(right.Output-left.Output)
		}
	}
	return points[len(points)-1].Output
}
