package main

import (
	"drum"
	"flag"
	"fmt"
	"log"
)

var fixturePath string

func main() {
	flag.StringVar(&fixturePath, "file", "../patterns/pattern_1.splice", "Path to a pattern (.splice) file")
	flag.Parse()
	p, err := drum.DecodeFile(fixturePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(p)
}
