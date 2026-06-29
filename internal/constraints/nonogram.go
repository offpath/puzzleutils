package constraints

import "github.com/offpath/puzzleutils/internal/puzzle"

type NonogramConstraint struct {
	variables []*puzzle.Variable
	lengths   []int
	p         *puzzle.Puzzle
}

func NewNonogramConstraint(p *puzzle.Puzzle, variables []*puzzle.Variable, lengths []int) *NonogramConstraint {
	return &NonogramConstraint{
		variables: variables,
		lengths:   lengths,
		p:         p,
	}
}

func (c *NonogramConstraint) Variables() []*puzzle.Variable {
	return c.variables
}

func (c *NonogramConstraint) Check() bool {
	b := NewBuildupSet(len(c.variables))
	var f func(lengths []int, ds []*puzzle.Variable)
	f = func(lengths []int, ds []*puzzle.Variable) {
		if len(ds) == 0 {
			return
		}
		sum := len(lengths) - 1
		for _, n := range lengths {
			sum += n
		}
		if sum > len(ds) {
			return
		}

		falseVal := c.p.GetBoolValue(false)
		trueVal := c.p.GetBoolValue(true)

		if ds[0].Values()[falseVal] {
			b.Push(falseVal)
			f(lengths, ds[1:])
			b.Pop()
		}
		if len(lengths) == 0 {
			return
		}
		for i := 0; i < lengths[0]; i++ {
			if !ds[i].Values()[trueVal] {
				return
			}
			b.Push(trueVal)
			defer b.Pop()
		}
		if len(ds) > lengths[0] {
			if !ds[lengths[0]].Values()[falseVal] {
				return
			}
			b.Push(falseVal)
			defer b.Pop()
			f(lengths[1:], ds[lengths[0]+1:])
		}
	}
	f(c.lengths, c.variables)
	b.Export(c.variables)

	if len(c.variables) > 0 && len(b.possibleSets[0]) == 0 {
		return false
	}
	return true
}
