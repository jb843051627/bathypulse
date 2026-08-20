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

func (t Threshold) Clamp(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > t.Value*2 {
		return t.Value * 2
	}
	return value
}

func (t Threshold) Label(value float64) string {
	if t.Rising(value, 0) {
		return "high"
	}
	return "normal"
}
