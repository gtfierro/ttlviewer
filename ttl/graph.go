package ttl

import (
	"github.com/d4l3k/turtle"
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
		if trimName(triple.Pred) == "type" {
			if trimName(triple.Obj) == "NamedIndividual" {
				continue
			}
			g.types[trimName(triple.Subj)] = trimName(triple.Obj)
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
	subj := trimName(triple.Subj)
	subType := g.types[subj]
	pred := trimName(triple.Pred)
	obj := trimName(triple.Obj)
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
