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

	for _, input := range tData {
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

	}
}
