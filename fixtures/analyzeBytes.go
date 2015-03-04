package main

// TODO: This program was coded in a sprint mostly while commuting on the L train, clean. it. up.

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var spliceSuffixRe = regexp.MustCompile("^.*\\.splice$")

func getSpliceFileInfos(path string) (spliceFileInfos []os.FileInfo) {
	fileInfos, err := ioutil.ReadDir(path)
	check(err)
	for _, fileInfo := range fileInfos {
		if spliceSuffixRe.Match([]byte(fileInfo.Name())) {
			spliceFileInfos = append(spliceFileInfos, fileInfo)
		}
	}
	return
}

func readFiles(fileInfos []os.FileInfo) map[string][]byte {
	allFiles := make(map[string][]byte, 0)
	for _, fileInfo := range fileInfos {
		fileContents, err := ioutil.ReadFile(fileInfo.Name())
		check(err)
		allFiles[fileInfo.Name()] = fileContents
	}
	return allFiles
}

func getLongestFileLengthInBytes(files map[string][]byte) (longest int, allEqual bool) {
	allEqual = true
	longest = -1
	for _, contents := range files {
		length := len(contents)
		if length > longest {
			longest = length
			allEqual = false
		}
	}
	return
}

func getMapKeys(aMap map[string][]byte) (keys []string) {
	for key, _ := range aMap {
		keys = append(keys, key)
	}
	return
}

type byteDelta struct { // TODO: A less horrible no good very bad name.
	uniform    bool
	valueFreqs map[byte]int
}

func main() {
	path := "/home/tasm/go/src/splice/fixtures" // TODO: No hardcoding.
	fileInfos := getSpliceFileInfos(path)
	allFiles := readFiles(fileInfos)
	longest, _ := getLongestFileLengthInBytes(allFiles)
	byteDeltas := make([]byteDelta, longest)
	fileNames := getMapKeys(allFiles)
	for i := 0; i < longest; i++ {
		byteDelta := byteDelta{uniform: false, valueFreqs: make(map[byte]int)}
		checkedAllFiles := true
		for _, fileName := range fileNames {
			if len(allFiles[fileName]) > i {
				byteAtOffset := allFiles[fileName][i]
				byteDelta.valueFreqs[byteAtOffset] += 1
			} else {
				checkedAllFiles = false
			}

		}
		if len(byteDelta.valueFreqs) == 1 && checkedAllFiles {
			byteDelta.uniform = true
		}
		byteDeltas[i] = byteDelta
	}
	for i, byteDelta := range byteDeltas {
		keys := make([]string, len(byteDelta.valueFreqs))
		for key, _ := range byteDelta.valueFreqs {
			keys = append(keys, string(key))
		}
		fmt.Println(i, byteDelta.uniform, byteDelta.valueFreqs, keys)
	}
}
