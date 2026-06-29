package layouts

import "github.com/offpath/puzzleutils/internal/puzzle"

type Node struct {
	ID   int
	Arcs []*Arc
}

func (n *Node) ArcTo(other *Node) *Arc {
	for _, a := range n.Arcs {
		if a.Node1 == other || a.Node2 == other {
			return a
		}
	}
	return nil
}

type Arc struct {
	ID    int
	Node1 *Node
	Node2 *Node
	Var   *puzzle.Variable
}

type Graph struct {
	p     *puzzle.Puzzle
	nodes []*Node
	arcs  []*Arc
}

func NewGraph(p *puzzle.Puzzle) *Graph {
	return &Graph{
		p: p,
	}
}

func (g *Graph) AddNode() *Node {
	n := &Node{
		ID: len(g.nodes),
	}
	g.nodes = append(g.nodes, n)
	return n
}

func (g *Graph) AddArc(n1, n2 *Node) *Arc {
	v := g.p.NewVariable()
	v.SetValueRange(puzzle.ValueSet{
		g.p.GetBoolValue(true):  true,
		g.p.GetBoolValue(false): true,
	})

	a := &Arc{
		ID:    len(g.arcs),
		Node1: n1,
		Node2: n2,
		Var:   v,
	}
	g.arcs = append(g.arcs, a)
	n1.Arcs = append(n1.Arcs, a)
	n2.Arcs = append(n2.Arcs, a)
	return a
}

func (g *Graph) Nodes() []*Node {
	return g.nodes
}

func (g *Graph) Arcs() []*Arc {
	return g.arcs
}

type SquareGridGraph struct {
	graph      *Graph
	rows, cols int
	nodes      [][]*Node
}

func NewSquareGridGraph(p *puzzle.Puzzle, rows, cols int) *SquareGridGraph {
	g := &SquareGridGraph{
		graph: NewGraph(p),
		rows:  rows,
		cols:  cols,
		nodes: make([][]*Node, rows),
	}

	for r := 0; r < rows; r++ {
		g.nodes[r] = make([]*Node, cols)
		for c := 0; c < cols; c++ {
			g.nodes[r][c] = g.graph.AddNode()
		}
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if c+1 < cols {
				g.graph.AddArc(g.nodes[r][c], g.nodes[r][c+1])
			}
			if r+1 < rows {
				g.graph.AddArc(g.nodes[r][c], g.nodes[r+1][c])
			}
		}
	}

	return g
}

func (g *SquareGridGraph) Get(row, col int) *Node {
	return g.nodes[row][col]
}

func (g *SquareGridGraph) Graph() *Graph {
	return g.graph
}

func (g *SquareGridGraph) SurroundingArcs(row, col int) []*Arc {
	n00 := g.Get(row, col)
	n01 := g.Get(row, col+1)
	n10 := g.Get(row+1, col)
	n11 := g.Get(row+1, col+1)

	var result []*Arc
	if a := n00.ArcTo(n01); a != nil {
		result = append(result, a)
	}
	if a := n00.ArcTo(n10); a != nil {
		result = append(result, a)
	}
	if a := n01.ArcTo(n11); a != nil {
		result = append(result, a)
	}
	if a := n10.ArcTo(n11); a != nil {
		result = append(result, a)
	}
	return result
}
