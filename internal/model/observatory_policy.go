package model

import "time"

type ObservatoryPolicy struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	RatioThreshold float64       `json:"ratio_threshold"`
	MaxWindow      time.Duration `json:"max_window"`
	ReviewAfter    time.Duration `json:"review_after"`
	Enabled        bool          `json:"enabled"`
}

func (p ObservatoryPolicy) Valid() bool {
	return p.ID != "" && p.Name != "" && p.RatioThreshold > 0 && p.MaxWindow > 0 && p.ReviewAfter > 0
}

func (p ObservatoryPolicy) ShouldReview(age time.Duration) bool {
	return p.Enabled && age >= p.ReviewAfter
}

func (p ObservatoryPolicy) StrongRatio(peak, rms float64) bool {
	return rms > 0 && peak/rms >= p.RatioThreshold
}

func (p ObservatoryPolicy) WindowAllowed(duration time.Duration) bool {
	return duration > 0 && duration <= p.MaxWindow
}
