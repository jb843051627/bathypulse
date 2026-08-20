package validation

import (
	"fmt"
	"strings"

	"github.com/jb843051627/bathypulse/internal/model"
)

func Packet(packet model.TelemetryPacket) error {
	if packet.ID == "" || packet.StationID == "" || packet.Type == "" {
		return fmt.Errorf("packet identity: %w", model.ErrInvalid)
	}
	if packet.Sequence <= 0 || len(packet.Payload) == 0 {
		return fmt.Errorf("packet body: %w", model.ErrInvalid)
	}
	if strings.TrimSpace(packet.Signature) == "" {
		return fmt.Errorf("packet signature: %w", model.ErrInvalid)
	}
	return nil
}

func IsStrongWaveform(peak, rms float64) bool {
	if rms <= 0 || peak <= rms {
		return false
	}
	return peak/rms >= 3.5
}
