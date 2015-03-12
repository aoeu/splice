package main

import (
	"splice/encoding/drum"
	"flag"
	"fmt"
	"log"
)

var fixturePath string
var defaultPath = "../encoding/drum/patterns/pattern_1.splice"

func main() {
	flag.StringVar(&fixturePath, "file", defaultPath, "Path to a pattern (.splice) file")
	flag.Parse()
	p, err := drum.DecodeFile(fixturePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(p)
}
