package drum

import (
	"bytes"
	"encoding/binary"
)

var unknownIndexes = map[int]struct{}{13: {}, 45: {}, 46: {}, 49: {}}

// TODO: Implement method that encodes from pattern wrap it.

// Encode a text file backup in a custom binary data format.
// TODO: Spec out the documentation.
func Encode(input string) ([]byte, error) {
	// TODO: Could this be bound to a Pattern type to be more consistent with Decode?
	p, err := NewPatternFromBackup(input)
	if err != nil {
		return []byte{}, err
	}
	// TODO: Implement header() method bound to pattern type.
	h := header{}
	h.ChunkID = [6]byte{'S', 'P', 'L', 'I', 'C', 'E'}
	// TODO: Smarter handling of tempos for variant versions.
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
		if err := binary.Write(buf, binary.LittleEndian, t.encode()); err != nil {
			return []byte{}, err
		}
	}
	return buf.Bytes(), nil
}
