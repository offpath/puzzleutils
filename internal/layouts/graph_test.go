package layouts

import (
	"testing"

	"github.com/offpath/puzzleutils/internal/puzzle"
)

func TestGraph(t *testing.T) {
	p := puzzle.NewPuzzle2()
	g := NewGraph(p)

	n1 := g.AddNode()
	n2 := g.AddNode()
	n3 := g.AddNode()

	a1 := g.AddArc(n1, n2)
	a2 := g.AddArc(n2, n3)

	if len(g.Nodes()) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(g.Nodes()))
	}
	if len(g.Arcs()) != 2 {
		t.Errorf("Expected 2 arcs, got %d", len(g.Arcs()))
	}

	if len(n1.Arcs) != 1 || n1.Arcs[0] != a1 {
		t.Errorf("Node 1 arcs incorrect")
	}
	if len(n2.Arcs) != 2 || n2.Arcs[0] != a1 || n2.Arcs[1] != a2 {
		t.Errorf("Node 2 arcs incorrect")
	}
	if len(n3.Arcs) != 1 || n3.Arcs[0] != a2 {
		t.Errorf("Node 3 arcs incorrect")
	}

	if len(a1.Var.Values()) != 2 {
		t.Errorf("Expected arc variable to have 2 possible values")
	}
}
