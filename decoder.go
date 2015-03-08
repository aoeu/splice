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
	h := header{}
	p := NewPattern()
	f, err := os.Open(path)
	if err != nil {
		return p, err
	}
	defer f.Close()
	if binary.Read(f, binary.LittleEndian, &h); err != nil {
		return p, err
	}
	p.HardwareVersion = string(bytes.Trim(h.HardwareVersion[:], "\x00"))
	reader := io.Reader(f)
	switch p.HardwareVersion {
	case "0.808-alpha":
		p.Tempo = int(h.Tempo) / 2
		if h.TempoDecimal != 0 {
			// TODO(aoeu): Is this really the correct way to determine the decimal?
			p.TempoDecimal = int(h.TempoDecimal) - 200
		}
	case "0.909":
		// TODO(aoeu): Is there no byte this value can be pulled from?
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

func (p Pattern) String() string {
	bpm := fmt.Sprint(p.Tempo)
	if p.TempoDecimal != 0 {
		// TODO(aoeu): Is this really the correct way to determine the decimal?
		bpm = fmt.Sprintf("%v.%v", p.Tempo, p.TempoDecimal)
	}
	s := fmt.Sprintf("Saved with HW Version: %s\nTempo: %v\n%v", p.HardwareVersion, bpm, p.Tracks)
	return s
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
