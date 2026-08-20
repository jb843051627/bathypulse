package model

type Threshold struct {
	Name       string  `json:"name"`
	Value      float64 `json:"value"`
	Hysteresis float64 `json:"hysteresis"`
}

func (t Threshold) Valid() bool {
	return t.Name != "" && t.Value > 0 && t.Hysteresis >= 0 && t.Hysteresis < t.Value
}

func (t Threshold) Rising(value, previous float64) bool {
	return value >= t.Value && previous < t.Value-t.Hysteresis
}

func (t Threshold) Falling(value, previous float64) bool {
	return value < t.Value-t.Hysteresis && previous >= t.Value
}

func (t Threshold) Normalize(value float64) float64 {
	if value <= 0 {
		return 0
	}
	return value / t.Value
}
