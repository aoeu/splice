package main

import (
	"drum"
	"flag"
	"fmt"
	"log"
)

var fixturePath string
var pattern string = `Saved with HW Version: 0.808-alpha
Tempo: 120
(0) kick	|x---|x---|x---|x---|
(1) snare	|----|x---|----|x---|
(2) clap	|----|x-x-|----|----|
(3) hh-open	|--x-|--x-|x-x-|--x-|
(4) hh-close	|x---|x---|----|x--x|
(5) cowbell	|----|----|--x-|----|
`

func main() {
	flag.StringVar(&fixturePath, "file", "../fixtures/pattern_1.splice", "Path to a fixture (.splice) file")
	flag.Parse()
	p, err := drum.DecodeFile(fixturePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(p)
}
