// Package drum is supposed to implement the decoding of .splice drum machine files.
// See golang-challenge.com/go-challenge1/ for more information
package drum

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

func NewPatternFromBackup(s string) (*Pattern, error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	p := NewPattern()
	for lineNum := 0; scanner.Scan() && lineNum < 2; lineNum++ {
		line := scanner.Text()
		switch lineNum {
		case 0:
			p.HardwareVersion = parseHardwareVersion(line)
		case 1:
			var err error
			p.Tempo, p.TempoDecimal, err = parseTempo(line)
			if err != nil {
				return p, err
			}
		default:
			p.Tracks = append(p.Tracks, parseTrack(line))
		}
	}
	return p, nil
}

func parseHardwareVersion(line string) string {
	s := strings.TrimLeft(line, "Saved with HW Version: ")
	return s
}

func parseTempo(line string) (tempo, tempoDecimal int, err error) {
	s := strings.TrimLeft(line, "Tempo: ")
	tempoRe := regexp.MustCompile("(\\d+).?(\\d+)?")
	match := tempoRe.FindStringSubmatch(s)
	tempo, err = strconv.Atoi(match[1])
	if err != nil {
		return 0, 0, err
	}
	tempoDecimal, err = strconv.Atoi(match[2])
	if err != nil {
		return tempo, 0, nil
	}
	return tempo, tempoDecimal, nil
}

func parseTrack(line string) Track {
	// TODO(aoeu): Implement
	return Track{}
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
