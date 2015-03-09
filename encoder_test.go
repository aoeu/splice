package drum

import (
	"io/ioutil"
	"path"
	"testing"
)

func TestEncodeTrack(t *testing.T) {
	// TODO(aoeu): Implement.
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

	for _, input := range tData {
		actual, err := Encode(input.backup)
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
