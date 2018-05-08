package ttl

import (
	"github.com/gtfierro/hod/turtle"
	"strings"
)

type Graph struct {
	nodes map[string]node
	types map[string]string
}

func trimName(input string) string {
	s := strings.Split(input, "#")
	if len(s) > 1 {
		return s[1]
	}
	return input
}

func NewGraph(triples []turtle.Triple) *Graph {
	g := &Graph{
		nodes: make(map[string]node),
		types: make(map[string]string),
	}

	// figure out what type all of the names are, and
	// remove the "type" triples from the list
	var newTriples []turtle.Triple
	for _, triple := range triples {
		if triple.Predicate.Value == "type" {
			if triple.Object.Value == "NamedIndividual" {
				continue
			}
			g.types[triple.Subject.Value] = triple.Object.Value
		} else {
			newTriples = append(newTriples, triple)
		}
	}
	triples = newTriples

	// add the triples to the graph
	for _, triple := range triples {
		g.addTriple(triple)
	}
	return g
}

func (g *Graph) addTriple(triple turtle.Triple) {
	subj := triple.Subject.Value
	subType := g.types[subj]
	pred := triple.Predicate.Value
	obj := triple.Object.Value
	objType := g.types[obj]
	// make subject node
	if _, found := g.nodes[subType]; !found {
		g.nodes[subType] = node{
			name:  subj,
			otype: g.types[subj],
			edges: []Edge{},
		}
	}
	// make object node
	if _, found := g.nodes[objType]; !found {
		g.nodes[objType] = node{
			name:  obj,
			otype: g.types[obj],
			edges: []Edge{},
		}
	}

	// connect them
	sNode := g.nodes[subType]
	oNode := g.nodes[objType]
	for _, edge := range sNode.edges {
		if edge.pred == pred && edge.node.otype == oNode.otype {
			return
		}
	}
	sNode.edges = append(sNode.edges, Edge{pred: pred, node: oNode})
	g.nodes[subType] = sNode
}

type Edge struct {
	pred string
	node node
}

type node struct {
	name  string
	otype string
	edges []Edge
}
