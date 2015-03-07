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
	if binary.Read(f, binary.LittleEndian, &p.Header); err != nil {
		return p, err
	}
	for {
		drumPart, err := readDrumPart(f)
		if err != nil {
			if err == io.EOF {
				return p, nil
			}
			return p, err
		}
		fmt.Sprintln("%+v %v", drumPart, err)
		p.DrumParts = append(p.DrumParts, drumPart)
	}
	return p, nil
}

type DrumParts []DrumPart

func (d DrumParts) String() string {
	s := ""
	for _, drumPart := range d {
		s += fmt.Sprintf("%v\n", drumPart)
	}
	return s
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	Header
	DrumParts
}

func NewPattern() (p *Pattern) {
	p = new(Pattern)
	p.DrumParts = make([]DrumPart, 0)
	return
}

var i int = 1

func (p Pattern) String() string {
	h := bytes.Trim(p.Header.HardwareVersion[:], "\x00")
	i++
	bpm := fmt.Sprintf("%v", (p.Header.BPM / 2))
	if p.Header.BPMDecimal != 0 {
		// TODO: Is this really the correct way to determine the decimal?
		bpm = fmt.Sprintf("%v.%v", bpm, p.Header.BPMDecimal-200)
	}
	s := fmt.Sprintf("Saved with HW Version: %s\nTempo: %v\n%v", h, bpm, p.DrumParts)
	return s
}

type Header struct {
	ChunkID         [6]byte  // 0 - 5
	Padding1        [7]byte  // 6
	Unknown1        [1]byte  // 13
	HardwareVersion [31]byte // 14 - 45
	Unknown2	[2]byte
	BPMDecimal      byte
	BPM             byte
	Unknown3        byte
}

type DrumPart struct {
	ID       byte
	Name     string
	Sequence []byte
}

func (d DrumPart) String() string {
	separator := "|"
	onBeat := "x"
	offBeat := "-"
	s := fmt.Sprintf("(%d) %s\t", d.ID, d.Name)
	for i := 0; i < len(d.Sequence); i++ {
		if i%4 == 0 {
			s += separator
		}
		if d.Sequence[i] == 1 {
			s += onBeat
		} else {
			s += offBeat
		}
	}
	s += separator
	return s
}

func NewDrumPart() *DrumPart {
	d := new(DrumPart)
	d.Sequence = make([]byte, 16)
	return d
}

func readDrumPart(r io.Reader) (d DrumPart, err error) {
	d = *NewDrumPart()
	if err = binary.Read(r, binary.LittleEndian, &d.ID); err != nil {
		return
	}
	padding := make([]byte, 3)
	if err = binary.Read(r, binary.LittleEndian, &padding); err != nil {
		return
	}
	var nameLen byte
	if err = binary.Read(r, binary.LittleEndian, &nameLen); err != nil {
		return
	}
	nameBytes := make([]byte, nameLen)
	if err = binary.Read(r, binary.LittleEndian, &nameBytes); err != nil {
		return
	}
	d.Name = string(nameBytes)
	if err = binary.Read(r, binary.LittleEndian, &d.Sequence); err != nil {
		return
	}
	return
}
