package report

import (
	"sort"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
)

type PolicyReport struct {
	GeneratedAt time.Time               `json:"generated_at"`
	Policy      model.ObservatoryPolicy `json:"policy"`
	Waveforms   int                     `json:"waveforms"`
	Strong      int                     `json:"strong"`
	Rejected    int                     `json:"rejected"`
	Ratios      []float64               `json:"ratios"`
}

func BuildPolicyReport(at time.Time, policy model.ObservatoryPolicy, waveforms []model.Waveform) PolicyReport {
	result := PolicyReport{GeneratedAt: at, Policy: policy, Waveforms: len(waveforms), Ratios: make([]float64, 0, len(waveforms))}
	for _, waveform := range waveforms {
		if waveform.RMS <= 0 {
			result.Rejected++
			continue
		}
		ratio := waveform.Peak / waveform.RMS
		result.Ratios = append(result.Ratios, ratio)
		if policy.StrongRatio(waveform.Peak, waveform.RMS) {
			result.Strong++
		}
	}
	sort.Float64s(result.Ratios)
	return result
}

func (p PolicyReport) StrongRatio() float64 {
	if p.Waveforms == 0 {
		return 0
	}
	return float64(p.Strong) / float64(p.Waveforms)
}
