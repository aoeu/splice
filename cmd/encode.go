package main

import (
	"drum"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

var patternPath string

func main() {
	flag.StringVar(&patternPath, "file",
		"../patterns/decoded/pattern_3.txt",
		"Path to a text file representing a pattern.")
	flag.Parse()
	data, err := ioutil.ReadFile(patternPath)
	if err != nil {
		log.Fatal(err)
	}
	b, err := drum.EncodeFile(string(data))
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Stdout.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}
