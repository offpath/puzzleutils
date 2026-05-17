package constraints2

import "github.com/offpath/puzzleutils/internal/puzzle"

type SetCountConstraint struct {
	variables  []*puzzle.Variable
	s          map[*puzzle.Value]int
	isCovering bool
}

func NewSetCountConstraint(variables []*puzzle.Variable, s map[*puzzle.Value]int, isCovering bool) *SetCountConstraint {
	return &SetCountConstraint{
		variables:  variables,
		s:          s,
		isCovering: isCovering,
	}
}

func (c *SetCountConstraint) Variables() []*puzzle.Variable {
	return c.variables
}

func (c *SetCountConstraint) Check() bool {
	for item, target := range c.s {
		count := 0
		possibleCount := 0
		for _, v := range c.variables {
			if v.Values()[item] {
				possibleCount++
				if len(v.Values()) == 1 {
					count++
				}
			}
		}
		if count > target || (c.isCovering && possibleCount < target) {
			return false
		}
		if count == target {
			for _, v := range c.variables {
				if len(v.Values()) > 1 || !v.Values()[item] {
					v.EliminateValue(item)
				}
			}
		}
	}
	return true
}
