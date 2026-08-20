package ingest

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"

	"github.com/jb843051627/bathypulse/internal/model"
)

type Frame struct {
	Version uint8
	Type    model.PacketType
	Seq     uint32
	Body    []byte
	CRC     uint32
}

func ParseFrame(data []byte) (Frame, error) {
	if len(data) < 9 {
		return Frame{}, fmt.Errorf("frame header: %w", model.ErrInvalid)
	}
	frame := Frame{Version: data[0], Type: model.PacketType(data[1]), Seq: binary.BigEndian.Uint32(data[2:6]), Body: model.CloneBytes(data[6 : len(data)-4]), CRC: binary.BigEndian.Uint32(data[len(data)-4:])}
	if crc32.ChecksumIEEE(frame.Body) != frame.CRC {
		return Frame{}, model.ErrChecksum
	}
	return frame, nil
}

func EncodeFrame(frame Frame) []byte {
	data := make([]byte, 10+len(frame.Body))
	data[0] = frame.Version
	data[1] = byte(len(frame.Type))
	binary.BigEndian.PutUint32(data[2:6], frame.Seq)
	copy(data[6:], frame.Body)
	binary.BigEndian.PutUint32(data[len(data)-4:], crc32.ChecksumIEEE(frame.Body))
	return data
}
