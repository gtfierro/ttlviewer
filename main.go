package main

import (
	"flag"
	"fmt"
	"github.com/gtfierro/ttlviewer/ttl"
	"log"
	"os"
)

var outputFile = flag.String("o", "output_graph.pdf", "Output file name for compiled graph")
var keepDot = flag.Bool("k", false, "Whether or not to keep generated DOT file")

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage:\nttlviewer <ttl file>\nttlviewer -o <output pdf> -k {false|true} <ttl file>")
		os.Exit(1)
	}

	reader, err := os.Open(flag.Args()[flag.NArg()-1])
	if err != nil {
		log.Panic(err)
	}
	output, _, err := ttl.RunFile(reader, *keepDot)
	if err != nil {
		log.Panic(err)
	}
	f, err := os.Create(*outputFile)
	if err != nil {
		log.Panic(err)
	}

	_, err = f.Write(output)

	if err != nil {
		log.Panic(err)
	}
}
