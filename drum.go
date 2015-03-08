// Package drum is supposed to implement the decoding of .splice drum machine files.
// See golang-challenge.com/go-challenge1/ for more information
package drum

import (
	"fmt"
)

type header struct {
	// TODO(aoeu): What are better names for "Unknown.*"?
	ChunkID         [6]byte  // 0 - 5
	Padding1        [7]byte  // 6
	Unknown1        [1]byte  // 13
	HardwareVersion [31]byte // 14 - 45
	Unknown2        [2]byte
	TempoDecimal    byte // Tempo Decimal for 808
	Tempo           byte // Tempo for 808
	Unknown3        byte
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	Tempo           int
	TempoDecimal    int
	HardwareVersion string
	Tracks
}

// NewPattern returns an empty pattern.
func NewPattern() *Pattern {
	// TODO(aoeu): json.NewDecoder(r io.Reader) could be influential.
	p := new(Pattern)
	p.Tracks = make([]Track, 0)
	return p
}

func NewPatternFromBackup(input string) (Pattern, error) {
	return Pattern{}, nil
}

// A Track represents a named, identified drum sequence.
type Track struct {
	ID       byte
	Name     string
	Sequence []byte
}

// NewTrack returns an empty, initialized track.
func NewTrack() *Track {
	t := new(Track)
	t.Sequence = make([]byte, 16)
	return t
}

const (
	separator string = "|"
	onBeat    string = "x"
	offBeat   string = "-"
	errorRune string = "?"
)

func (t Track) String() string {
	s := fmt.Sprintf("(%d) %s\t", t.ID, t.Name)
	for i := 0; i < len(t.Sequence); i++ {
		if i%4 == 0 {
			s += separator
		}
		switch t.Sequence[i] {
		case 1:
			s += onBeat
		case 0:
			s += offBeat
		default:
			s += errorRune
		}
	}
	s += separator
	return s
}
