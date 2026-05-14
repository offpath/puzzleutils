package constraints2

import "github.com/offpath/puzzleutils/internal/puzzle"

type BuildupSet struct {
	size, cursor int
	values       []*puzzle.Value
	possibleSets []puzzle.ValueSet
}

func NewBuildupSet(size int) *BuildupSet {
	result := &BuildupSet{
		size:         size,
		cursor:       0,
		values:       make([]*puzzle.Value, size),
		possibleSets: make([]puzzle.ValueSet, size),
	}
	for i := 0; i < size; i++ {
		result.possibleSets[i] = puzzle.ValueSet{}
	}
	return result
}

func (b *BuildupSet) Push(val *puzzle.Value) {
	b.values[b.cursor] = val
	b.cursor++
	if b.cursor == b.size {
		for i := 0; i < b.size; i++ {
			b.possibleSets[i][b.values[i]] = true
		}
	}
}

func (b *BuildupSet) Pop() {
	b.cursor--
}

func (b *BuildupSet) Export(variables []*puzzle.Variable) {
	for i, s := range b.possibleSets {
		for val := range variables[i].Values() {
			if !s[val] {
				variables[i].EliminateValue(val)
			}
		}
	}
}
