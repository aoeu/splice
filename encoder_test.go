package drum

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"
)

func TestEncodeTrack(t *testing.T) {
	// TODO(aoeu): Implement.
	expected := []byte{2, 0, 0, 0,
		4, 'c', 'l', 'a', 'p',
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 1, 0, 1,
		0, 0, 1, 1}
	track := Track{ID: 2, Name: "clap",
		Sequence: []byte{1, 0, 0, 0,
			0, 1, 0, 0,
			0, 1, 0, 1,
			0, 0, 1, 1}}
	actual := track.encode()
	if len(expected) != len(actual) {
		t.Fatalf("Expected %v output bytes and got %v", len(expected), len(actual))
	}
	for i, b := range actual {
		if expected[i] != b {
			t.Fatalf("Expected '%v' byte but received '%v' at %v", expected[i], b, i)
		}
	}
}

func TestEncode(t *testing.T) {
	tData := []struct {
		path   string
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

	b := new(bytes.Buffer)
	e := NewEncoder(b)
	for _, input := range tData {
		pattern, err := NewPatternFromBackup(input.backup)
		err = e.Encode(*pattern)
		actual := b.Bytes()
		if err != nil {
			t.Fatalf("Something went wrong encoding - %v", err)
		}
		filePath := path.Join("patterns", input.path+".splice")
		expected, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Could not read input file at path %v - %v", filePath, err)
		}
		if len(expected) != len(actual) {
			t.Fatalf("Expected %v output bytes and got %v", len(expected), len(actual))
		}
		for i, b := range actual {
			if _, ok := unknownIndexes[i]; ok {
				continue
			}
			if expected[i] != b {
				t.Fatalf("Expected '%v' byte but received '%v' at %v", expected[i], b, i)
			}
		}
	}
}
