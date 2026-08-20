package model

type QualityFlag string

const (
	QualityGood        QualityFlag = "good"
	QualityNeedsReview QualityFlag = "review"
	QualityDiscarded   QualityFlag = "discarded"
	QualityClockSkew   QualityFlag = "clock_skew"
	QualitySaturated   QualityFlag = "saturated"
)

type QualityReview struct {
	WaveformID string      `json:"waveform_id"`
	Flag       QualityFlag `json:"flag"`
	Reason     string      `json:"reason"`
	Score      float64     `json:"score"`
}

func (q QualityReview) Accepted() bool {
	return q.Flag == QualityGood || (q.Flag == QualityNeedsReview && q.Score >= 0.75)
}

func (q QualityReview) NeedsOperator() bool {
	return q.Flag == QualityNeedsReview || q.Flag == QualityClockSkew
}
