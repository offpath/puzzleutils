package constraints2

import (
	"github.com/offpath/puzzleutils/internal/layouts"
	"github.com/offpath/puzzleutils/internal/puzzle"
)

type GridGraphPointConstraint struct {
	node *layouts.Node
	p    *puzzle.Puzzle
}

func NewGridGraphPointConstraint(p *puzzle.Puzzle, node *layouts.Node) *GridGraphPointConstraint {
	return &GridGraphPointConstraint{node: node, p: p}
}

func (c *GridGraphPointConstraint) Variables() []*puzzle.Variable {
	var vars []*puzzle.Variable
	for _, arc := range c.node.Arcs {
		vars = append(vars, arc.Var)
	}
	return vars
}

func (c *GridGraphPointConstraint) Check() bool {
	possible := 0
	set := 0
	trueVal := c.p.GetBoolValue(true)
	falseVal := c.p.GetBoolValue(false)

	for _, arc := range c.node.Arcs {
		v := arc.Var
		if len(v.Values()) == 1 && v.Values()[trueVal] {
			set++
		}
		if v.Values()[trueVal] {
			possible++
		}
	}

	if set > 2 || (set == 1 && possible == 1) {
		return false
	}

	if set == 1 && possible == 2 {
		for _, arc := range c.node.Arcs {
			v := arc.Var
			if v.Values()[trueVal] {
				v.RestrictTo(trueVal)
			}
		}
	} else if possible < 2 {
		for _, arc := range c.node.Arcs {
			v := arc.Var
			v.RestrictTo(falseVal)
		}
	}
	return true
}

type GridGraphBoxConstraint struct {
	arcs []*layouts.Arc
	n    int
	p    *puzzle.Puzzle
}

func NewGridGraphBoxConstraint(p *puzzle.Puzzle, arcs []*layouts.Arc, n int) *GridGraphBoxConstraint {
	return &GridGraphBoxConstraint{arcs: arcs, n: n, p: p}
}

func (c *GridGraphBoxConstraint) Variables() []*puzzle.Variable {
	var vars []*puzzle.Variable
	for _, arc := range c.arcs {
		vars = append(vars, arc.Var)
	}
	return vars
}

func (c *GridGraphBoxConstraint) Check() bool {
	possible := 0
	set := 0
	trueVal := c.p.GetBoolValue(true)
	falseVal := c.p.GetBoolValue(false)

	for _, arc := range c.arcs {
		v := arc.Var
		if len(v.Values()) == 1 && v.Values()[trueVal] {
			set++
		}
		if v.Values()[trueVal] {
			possible++
		}
	}

	if possible < c.n || set > c.n {
		return false
	}
	if possible == c.n {
		for _, arc := range c.arcs {
			v := arc.Var
			if v.Values()[trueVal] {
				v.RestrictTo(trueVal)
			}
		}
		set = c.n
	}
	if set == c.n {
		for _, arc := range c.arcs {
			v := arc.Var
			if len(v.Values()) > 1 || !v.Values()[trueVal] {
				v.RestrictTo(falseVal)
			}
		}
	}
	return true
}

type GridGraphLoopConstraint struct {
	graph *layouts.Graph
	p     *puzzle.Puzzle
}

func NewGridGraphLoopConstraint(p *puzzle.Puzzle, graph *layouts.Graph) *GridGraphLoopConstraint {
	return &GridGraphLoopConstraint{graph: graph, p: p}
}

func (c *GridGraphLoopConstraint) Variables() []*puzzle.Variable {
	var vars []*puzzle.Variable
	for _, arc := range c.graph.Arcs() {
		vars = append(vars, arc.Var)
	}
	return vars
}

func (c *GridGraphLoopConstraint) Check() bool {
	hasLoop := false
	hasUnfinishedLoop := false
	trueVal := c.p.GetBoolValue(true)

	for _, startNode := range c.graph.Nodes() {
		curNode := startNode
		seenNodes := map[*layouts.Node]bool{}
		var prevArc *layouts.Arc

		for {
			var arcFound *layouts.Arc
			for _, arc := range curNode.Arcs {
				if arc != prevArc && len(arc.Var.Values()) == 1 && arc.Var.Values()[trueVal] {
					arcFound = arc
					break
				}
			}

			if arcFound == nil {
				if len(seenNodes) > 0 {
					hasUnfinishedLoop = true
				}
				break
			}

			prevArc = arcFound
			seenNodes[curNode] = true

			if arcFound.Node1 == curNode {
				curNode = arcFound.Node2
			} else {
				curNode = arcFound.Node1
			}

			if seenNodes[curNode] {
				hasLoop = true
				break
			}
		}
	}

	return !(hasLoop && hasUnfinishedLoop)
}
