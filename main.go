package main

import (
	"flag"
	"fmt"
	"github.com/gtfierro/ttlviewer/ttl"
	"log"
	"os"
)

var outputFile = flag.String("o", "output_graph.pdf", "Output file name for compiled graph")

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage:\nttlviewer <ttl file>\nttlviewer -o <output pdf> <ttl file>")
		os.Exit(1)
	}

	reader, err := os.Open(flag.Args()[flag.NArg()-1])
	if err != nil {
		log.Panic(err)
	}
	output, err := ttl.RunFile(reader)
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
