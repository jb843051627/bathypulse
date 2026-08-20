package model

import "time"

type PacketType string

const (
	PacketWaveform  PacketType = "waveform"
	PacketHeartbeat PacketType = "heartbeat"
	PacketControl   PacketType = "control"
	PacketSnapshot  PacketType = "snapshot"
)

type TelemetryPacket struct {
	ID         string     `json:"id"`
	StationID  string     `json:"station_id"`
	Type       PacketType `json:"type"`
	ReceivedAt time.Time  `json:"received_at"`
	Sequence   int64      `json:"sequence"`
	Payload    []byte     `json:"payload"`
	Signature  string     `json:"signature"`
}

func (p TelemetryPacket) PayloadSize() int {
	return len(p.Payload)
}

func (p TelemetryPacket) IsControl() bool {
	return p.Type == PacketControl || p.Type == PacketSnapshot
}

func (p TelemetryPacket) Age(at time.Time) time.Duration {
	return at.Sub(p.ReceivedAt)
}
