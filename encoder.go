package drum

import (
	"bytes"
	"encoding/binary"
)

var unknownIndexes = map[int]struct{}{13: {}, 45: {}, 46: {}, 49: {}}

func Encode(input string) ([]byte, error) {
	p, err := NewPatternFromBackup(input)
	if err != nil {
		return []byte{}, err
	}
	h := header{}
	h.ChunkID = [6]byte{'S', 'P', 'L', 'I', 'C', 'E'}
	for i, r := range p.HardwareVersion {
		h.HardwareVersion[i] = byte(r)
	}
	if p.TempoDecimal != 0 {
		h.TempoDecimal = byte(p.TempoDecimal + 200)
	}
	h.Tempo = byte(p.Tempo * 2)
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, h); err != nil {
		return []byte{}, nil
	}
	for _, t := range p.Tracks {
		if err := binary.Write(buf, binary.LittleEndian, t.Encode()); err != nil {
			return []byte{}, err
		}
	}
	return buf.Bytes(), nil
}
