package drum

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (*Pattern, error) {
	p := NewPattern()
	f, err := os.Open(path)
	if err != nil {
		return p, err
	}
	defer f.Close()
	if binary.Read(f, binary.LittleEndian, &p.header); err != nil {
		return p, err
	}
	p.HardwareVersion = string(bytes.Trim(p.header.HardwareVersion[:], "\x00"))
	reader := io.Reader(f)
	switch p.HardwareVersion {
	case "0.808-alpha":
		p.Tempo = int(p.header.Tempo) / 2
		if p.header.TempoDecimal != 0 {
			// TODO: Is this really the correct way to determine the decimal?
			p.TempoDecimal = int(p.header.TempoDecimal) - 200
		}
	case "0.909":
		// TODO: Is there no byte this value can be pulled from?
		p.Tempo = 240
	case "0.708-alpha":
		p.Tempo = 999
	}
	p.Tracks, err = readAllTracks(reader)
	if err == io.ErrUnexpectedEOF {
		return p, nil
	}
	return p, err
}

func readAllTracks(r io.Reader) (Tracks, error) {
	var d []Track
	for {
		drumPart, err := readTrack(r)
		if err != nil {
			if err == io.EOF {
				return d, nil
			}
			return d, err
		}
		d = append(d, drumPart)
	}
}

// Tracks is a drum Track series that comprises the pattern.
type Tracks []Track

func (d Tracks) String() string {
	s := ""
	for _, drumPart := range d {
		s += fmt.Sprintf("%v\n", drumPart)
	}
	return s
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	header
	Tempo           int
	TempoDecimal    int
	HardwareVersion string
	Tracks
}

// NewPattern returns an empty pattern.
func NewPattern() *Pattern {
	// TODO: json.NewDecoder(r io.Reader) could be influential.
	p := new(Pattern)
	p.Tracks = make([]Track, 0)
	return p
}

func (p Pattern) String() string {
	h := bytes.Trim(p.header.HardwareVersion[:], "\x00")
	bpm := fmt.Sprint(p.Tempo)
	if p.TempoDecimal != 0 {
		// TODO: Is this really the correct way to determine the decimal?
		bpm = fmt.Sprintf("%v.%v", p.Tempo, p.TempoDecimal)
	}
	s := fmt.Sprintf("Saved with HW Version: %s\nTempo: %v\n%v", h, bpm, p.Tracks)
	return s
}

type header struct {
	// TODO: What are better names for "Unknown.*"?
	ChunkID         [6]byte  // 0 - 5
	Padding1        [7]byte  // 6
	Unknown1        [1]byte  // 13
	HardwareVersion [31]byte // 14 - 45
	Unknown2        [2]byte
	TempoDecimal    byte // Tempo Decimal for 808
	Tempo           byte // Tempo for 808
	Unknown3        byte
}

// A Track represents a named, identified drum sequence.
type Track struct {
	ID       byte
	Name     string
	Sequence []byte
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

// NewTrack returns an empty, initialized track.
func NewTrack() *Track {
	t := new(Track)
	t.Sequence = make([]byte, 16)
	return t
}

func readTrack(r io.Reader) (Track, error) {
	t := *NewTrack()
	if err := binary.Read(r, binary.LittleEndian, &t.ID); err != nil {
		return t, err
	}
	padding := make([]byte, 3)
	if err := binary.Read(r, binary.LittleEndian, &padding); err != nil {
		return t, err
	}
	var nameLen byte
	if err := binary.Read(r, binary.LittleEndian, &nameLen); err != nil {
		return t, err
	}
	nameBytes := make([]byte, nameLen)
	if err := binary.Read(r, binary.LittleEndian, &nameBytes); err != nil {
		return t, err
	}
	t.Name = string(nameBytes)
	if err := binary.Read(r, binary.LittleEndian, &t.Sequence); err != nil {
		return t, err
	}
	return t, nil
}
