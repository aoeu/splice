package drum

import (
	"testing"
)

func TestParseHardwareVersion(t *testing.T) {
	expected := "0.808-alpha"
	actual := parseHardwareVersion("Saved with HW Version: " + expected)
	if actual != expected {
		t.Fatalf(`Expected "%v" but received "%v"`, expected, actual)
	}
}

func TestParseTempo(t *testing.T) {
	tData := []struct {
		input        string
		tempo        int
		tempoDecimal int
	}{
		{"Tempo: 120", 120, 0},
		{"Tempo: 99", 99, 0},
		{"Tempo: 91.3", 91, 3},
	}
	for _, expected := range tData {
		tempo, tempoDecimal, err := parseTempo(expected.input)
		if err != nil {
			t.Fatal(err)
		}
		if tempo != expected.tempo {
			t.Fatal("Expected tempo %v but received %v", expected.tempo, tempo)
		}
		if tempoDecimal != expected.tempoDecimal {
			t.Fatal("Expected tempo decimal %v but received %v", expected.tempoDecimal, tempoDecimal)
		}

	}
}

func TestParseTrackId(t *testing.T) {
	input := "(3) hh-open	|--x-|--x-|x-x-|--x-|"
	var expectedId byte = 3
	expectedLine := "hh-open	|--x-|--x-|x-x-|--x-|"
	actualId, actualLine, err := parseTrackID(input)
	if err != nil {
		t.Fatal(err)
	}
	if expectedId != actualId {
		t.Fatalf("Expected ID %v but received %v", expectedId, actualId)
	}
	if expectedLine != actualLine {
		t.Fatalf("Expected line '%v' but received '%v'", expectedLine, actualLine)
	}
}

func TestParseTrackName(t *testing.T) {
	input := "hh-open	|--x-|--x-|x-x-|--x-|"
	expectedName := "hh-open"
	expectedLine := "--x-|--x-|x-x-|--x-|"
	actualName, actualLine := parseTrackName(input)
	if expectedName != actualName {
		t.Fatalf("Expected name '%v' but received '%v'", expectedName, actualName)
	}
	if expectedLine != actualLine {
		t.Fatalf("Expected line '%v' but received '%v'", expectedLine, actualLine)
	}

}

// func parseBar(line string, numMeasures int) (bar []byte, subLine string)
func TestParseBar(t *testing.T) {
	input := "--x-|--x-|x-x-|--x-|"
	expected := []byte{0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 1, 0}
	actual, s := parseBar(input, 4)
	for i, b := range actual {
		if expected[i] != b {
			t.Fatal("Expected '%v' but received '%v' at %v",
				expected, actual, i)
		}
	}
	if s != "" {
		t.Fatal("Expected empty string but received '%v'", s)
	}
}

func TestParseBeats(t *testing.T) {
	measure := "--x-"
	expected := []byte{0, 0, 1, 0}
	actual := parseBeats(measure)
	for i, b := range actual {
		if expected[i] != b {
			t.Fatal("Expected '%v' but received '%v' at %v",
				expected, actual, i)
		}
	}
}

func TestNewPatternFromBackup(t *testing.T) {
	tData := []struct {
		name   string
		backup string
	}{
		{"pattern_1",
			`Saved with HW Version: 0.808-alpha
Tempo: 120
(0) kick	|x---|x---|x---|x---|
(1) snare	|----|x---|----|x---|
(2) clap	|----|x-x-|----|----|
(3) hh-open	|--x-|--x-|x-x-|--x-|
(4) hh-close	|x---|x---|----|x--x|
(5) cowbell	|----|----|--x-|----|
`,
		},
	}

	input := tData[0]
	p, err := NewPatternFromBackup(input.backup)
	if err != nil {
		t.Fatal("Could not create Pattern from backup - %v", err)
	}
	if p.HardwareVersion != "0.808-alpha" {
		t.Fatalf("wrong version - %v", p.HardwareVersion)
	}
	if p.Tempo != 120 {
		t.Fatalf("wrong tempo - %v", p.Tempo)
	}
	expectedNumTracks := 6
	actualNumTracks := len(p.Tracks)
	if actualNumTracks != expectedNumTracks {
		t.Fatalf("Expected %v tracks but received %v tracks", expectedNumTracks, actualNumTracks)
	}
}
