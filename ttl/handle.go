package ttl

import (
	"github.com/d4l3k/turtle"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func RunFile(filename string) error {
	reader, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	out, err := turtle.Parse(reader)
	if err != nil {
		return err
	}

	g := ttl.NewGraph(out)

	for _, node := range g.nodes {
		fmt.Printf("> Name: %v, Type %v\n", node.name, node.otype)
		for _, edge := range node.edges {
			fmt.Printf("\t> Pred: %v, End Type %v\n", edge.pred, edge.node.otype)
		}
	}

	f, err := os.Create("output.gv")
	if err != nil {
		return err
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
	fmt.Println("output to", *outputFile)
	cmd := exec.Command("dot", "-Tpdf", "output.gv", "-o", *outputFile)
	return cmd.Run()
}
