package puzzle

import (
	"testing"

	"github.com/offpath/puzzleutils/internal/csp"
	"github.com/offpath/puzzleutils/internal/decide"
)

type sliceTracker struct {
	solutions [][]int
	vars      []*Variable
}

func (s *sliceTracker) CaptureSolution(prob *csp.Problem) bool {
	var sol []int
	for _, v := range s.vars {
		val := v.Values().Value()
		if val != nil {
			sol = append(sol, val.Int())
		} else {
			sol = append(sol, -1)
		}
	}
	s.solutions = append(s.solutions, sol)
	return false // returning false means "keep searching"
}

type notEqualConstraint struct {
	v1, v2 *Variable
}

func (c *notEqualConstraint) Variables() []*Variable { return []*Variable{c.v1, c.v2} }
func (c *notEqualConstraint) Check() bool {
	val1 := c.v1.Values().Value()
	val2 := c.v2.Values().Value()
	if val1 != nil && val2 != nil {
		return val1.Int() != val2.Int()
	}
	return true
}

func TestMultipleSolutions(t *testing.T) {
	p := NewPuzzle()
	v1 := p.NewVariable()
	v2 := p.NewVariable()

	v1.SetValueRange(p.GetIntRange(1, 3))
	v2.SetValueRange(p.GetIntRange(1, 3))

	p.AddConstraint(&notEqualConstraint{v1, v2})

	tracker := &sliceTracker{vars: []*Variable{v1, v2}}
	settings := csp.Settings{
		SolutionTracker: tracker,
		Decider:         &decide.First{},
	}

	success := p.Solve(settings)
	if !success {
		t.Errorf("expected Solve to return true, but it returned false")
	}

	// 1!=2, 1!=3, 2!=1, 2!=3, 3!=1, 3!=2 -> 6 solutions
	if len(tracker.solutions) != 6 {
		t.Errorf("expected 6 solutions, got %d: %v", len(tracker.solutions), tracker.solutions)
	}
}
