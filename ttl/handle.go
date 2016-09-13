package ttl

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"github.com/d4l3k/turtle"
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

	name := gethash() + ".gv"
	f, err := os.Create(name)
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
	cmd := exec.Command("dot", "-Tpdf", name)
	res, err := cmd.Output()
	os.Remove(name)
	return res, err
}
