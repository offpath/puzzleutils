package constraints

import (
	"github.com/offpath/puzzleutils/internal/puzzle"
	"github.com/offpath/puzzleutils/internal/trie"
)

type ValidWordConstraint struct {
	p         *puzzle.Puzzle
	variables []*puzzle.Variable
	t         *trie.Trie
	valueSet  puzzle.ValueSet
}

func NewValidWordConstraint(p *puzzle.Puzzle, variables []*puzzle.Variable, t *trie.Trie, valueSet puzzle.ValueSet) *ValidWordConstraint {
	return &ValidWordConstraint{p, variables, t, valueSet}
}

func (c *ValidWordConstraint) Variables() []*puzzle.Variable { return c.variables }

func (c *ValidWordConstraint) Check() bool {
	b := NewBuildupSet(len(c.variables))
	var f func(vars []*puzzle.Variable, prefix string)
	f = func(vars []*puzzle.Variable, prefix string) {
		if len(vars) == 0 {
			return
		}
		v := vars[0]
		for valObj := range c.valueSet {
			valStr := valObj.Str()
			if !v.Values()[valObj] {
				continue
			}
			if (len(vars) == 1 && c.t.HasWord(prefix+valStr)) || (len(vars) > 1 && c.t.HasPrefix(prefix+valStr)) {
				b.Push(valObj)
				f(vars[1:], prefix+valStr)
				b.Pop()
			}
		}
	}
	f(c.variables, "")
	b.Export(c.variables)
	return true
}
