package ttl

import (
	"fmt"
	"github.com/d4l3k/turtle"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func RunFile(ttl io.Reader) ([]byte, error) {
	var ret []byte
	// I KNOW not to use this; dirty hack until I have a better solution. Issue is the ttl parser only works on []byte
	bytes, err := ioutil.ReadAll(ttl)
	if err != nil {
		return ret, err
	}
	out, err := turtle.Parse(bytes)
	if err != nil {
		return ret, err
	}

	g := NewGraph(out)

	for _, node := range g.nodes {
		fmt.Printf("> Name: %v, Type %v\n", node.name, node.otype)
		for _, edge := range node.edges {
			fmt.Printf("\t> Pred: %v, End Type %v\n", edge.pred, edge.node.otype)
		}
	}

	f, err := os.Create("output.gv")
	if err != nil {
		return ret, err
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
	cmd := exec.Command("dot", "-Tpdf", "output.gv")
	return cmd.Output()
}
