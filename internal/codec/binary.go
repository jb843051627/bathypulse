package codec

import (
	"encoding/binary"
	"fmt"
	"math"
)

func DecodeFloat64s(data []byte) ([]float64, error) {
	if len(data)%8 != 0 {
		return nil, fmt.Errorf("float frame has %d bytes", len(data))
	}
	values := make([]float64, len(data)/8)
	for i := range values {
		bits := binary.BigEndian.Uint64(data[i*8:])
		values[i] = math.Float64frombits(bits)
	}
	return values, nil
}

func EncodeFloat64s(values []float64) []byte {
	data := make([]byte, len(values)*8)
	for i, value := range values {
		binary.BigEndian.PutUint64(data[i*8:], math.Float64bits(value))
	}
	return data
}
