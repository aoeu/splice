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
	Unknown2        [2]byte  //
	TempoDecimal    byte     // Tempo Decimal for 808
	Tempo           byte     // Tempo for 808
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

func (p Pattern) String() string {
	bpm := fmt.Sprint(p.Tempo)
	if p.TempoDecimal != 0 {
		bpm = fmt.Sprintf("%v.%v", p.Tempo, p.TempoDecimal)
	}
	s := fmt.Sprintf("Saved with HW Version: %s\nTempo: %v\n%v", p.HardwareVersion, bpm, p.Tracks)
	return s
}

// NewPatternFromBackup creates a pattern structure by parsing
// a backup file's human-readible text data.
func NewPatternFromBackup(s string) (*Pattern, error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	p := NewPattern()
	for lineNum := 0; scanner.Scan(); lineNum++ {
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
			t, err := parseTrack(line)
			if err != nil {
				return p, err
			}
			p.Tracks = append(p.Tracks, t)
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

var idRe = regexp.MustCompile(`\((\d+)\) `)
var beatRe = regexp.MustCompile(`([x-]{4})\|`)

func parseTrack(line string) (Track, error) {
	// TODO(aoeu): rename "Id" to "ID"
	id, line, err := parseTrackID(line)
	if err != nil {
		return Track{}, err
	}
	name, line := parseTrackName(line)
	bars, line := parseBar(line, 4)
	return Track{Name: name, ID: id, Sequence: bars}, nil
}

func parseTrackID(line string) (id byte, subLine string, err error) {
	idMatch := idRe.FindAllStringSubmatch(line, 1)[0]
	n, err := strconv.Atoi(idMatch[1])
	if err != nil {
		return id, subLine, err
	}
	subLine = strings.TrimLeft(line, idMatch[0])
	return byte(n), subLine, nil
}

var nameRe = regexp.MustCompile(`([\w-]+)\s+\|`)

func parseTrackName(line string) (name, subLine string) {
	s := strings.SplitN(line, "|", 2)
	name = strings.TrimRight(s[0], " \t")
	return name, s[1]
}

func parseBar(line string, numMeasures int) (bar []byte, subLine string) {
	measureMatch := beatRe.FindAllStringSubmatch(line, numMeasures)
	for i := 0; i < numMeasures; i++ {
		measure := measureMatch[i][1]
		bar = append(bar, parseBeats(measure)...)
		line = strings.TrimLeft(line, measureMatch[i][0])
	}
	return bar, line
}

func parseBeats(measure string) []byte {
	var beats []byte
	for _, beat := range measure {
		switch string(beat) {
		case onBeat:
			beats = append(beats, 1)
		case offBeat:
			beats = append(beats, 0)
		default:
			// TODO(aoeu): This doesn't seem correct.
			beats = append(beats, 0)
		}
	}
	return beats
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

func (t Track) encode() []byte {
	b := []byte{t.ID, 0, 0, 0}
	b = append(b, byte(len(t.Name)))
	b = append(b, []byte(t.Name)...)
	b = append(b, t.Sequence...)
	return b
}

const (
	separator string = "|"
	// TODO(aoeu): Could these be runes instead of strings?
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
