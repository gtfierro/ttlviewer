package main

import (
	"flag"
	"fmt"
	"github.com/d4l3k/turtle"
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

	err := ttl.RunFile(flag.Args()[flag.NArg()-1])
	if err != nil {
		log.Panic(err)
	}
}
