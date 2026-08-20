package model

import "time"

func (w Waveform) EndAt() time.Time {
	return w.CapturedAt.Add(time.Duration(w.DurationMs) * time.Millisecond)
}

func (w Waveform) Energy() float64 {
	return w.RMS * w.RMS * float64(w.DurationMs) / 1000
}

func (w Waveform) IsUsable() bool {
	return w.State != WaveformRejected && w.SampleRate > 0 && len(w.Payload) > 0
}
