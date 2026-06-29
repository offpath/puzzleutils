package constraints

import "github.com/offpath/puzzleutils/internal/puzzle"

type UniqueConstraint struct {
	variables []*puzzle.Variable
	vs        puzzle.ValueSet
}

func NewUniqueConstraint(variables []*puzzle.Variable, vs puzzle.ValueSet) *UniqueConstraint {
	return &UniqueConstraint{variables, vs}
}
func (c *UniqueConstraint) Variables() []*puzzle.Variable { return c.variables }
func (c *UniqueConstraint) Check() bool {
	for _, v1 := range c.variables {
		if val := v1.Values().Value(); val != nil {
			for _, v2 := range c.variables {
				if v1 != v2 {
					v2.EliminateValue(val)
				}
			}
		}
	}
	missingCount := 0
	for val := range c.vs {
		found := false
		for _, v := range c.variables {
			if v.Values()[val] {
				found = true
				break
			}
		}
		if !found {
			missingCount++
		}
	}
	return missingCount <= len(c.vs)-len(c.variables)
}
