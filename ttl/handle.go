package ttl

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"github.com/gtfierro/hod/turtle"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

func gethash() string {
	h := md5.New()
	seed := make([]byte, 16)
	binary.PutVarint(seed, time.Now().UnixNano())
	h.Write(seed)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// if keepDOT is True, then keep the DOT file around
func RunFile(ttl io.Reader, keepDOT bool) (pdf []byte, dot []byte, err error) {
	var ret []byte
	p := turtle.GetParser()
	dataset, _, err := p.ParseReader(ttl)
	if err != nil {
		return ret, ret, err
	}

	g := NewGraph(dataset.Triples)

	for _, node := range g.nodes {
		fmt.Printf("> Name: %v, Type %v\n", node.name, node.otype)
		for _, edge := range node.edges {
			fmt.Printf("\t> Pred: %v, End Type %v\n", edge.pred, edge.node.otype)
		}
	}

	name := gethash() + ".gv"
	f, err := os.Create(name)
	if err != nil {
		return ret, ret, err
	}

	fmt.Fprintln(f, "digraph G {")
	fmt.Fprintln(f, "ratio=\"auto\"")
	fmt.Fprintln(f, "rankdir=\"LR\"")
	fmt.Fprintln(f, "size=\"7.5,10\"")
	for _, node := range g.nodes {
		if node.otype == "" {
			continue
		}
		fmt.Fprintf(f, "%s;\n", node.otype)
		for _, edge := range node.edges {
			if edge.node.otype == "" {
				continue
			}
			fmt.Fprintf(f, "%s -> %s [label=\"%s\"];\n", node.otype, edge.node.otype, edge.pred)
		}
	}
	fmt.Fprintln(f, "}")
	cmd := exec.Command("dot", "-Tpdf", name)
	pdf, err = cmd.Output()
	if err != nil {
		return ret, ret, err
	}

	// grab contents of the DOT file
	file, err := os.Open(name)
	if err != nil {
		return ret, ret, err
	}
	dot, err = ioutil.ReadAll(file)
	if err != nil {
		return ret, ret, err
	}

	if !keepDOT {
		os.Remove(name)
	}
	return pdf, dot, err
}
