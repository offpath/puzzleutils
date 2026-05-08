package puzzle

import (
	"github.com/offpath/puzzleutils/internal/csp"
)

type Puzzle struct {
	problem  *csp.Problem
	valueSet []string
}

func NewPuzzle(size int, valueSet []string) *Puzzle {
	return &Puzzle{csp.NewProblem(size, len(valueSet)), valueSet}
}

func (p *Puzzle) AllGroup() []int {
	var result []int
	for i := 0; i < p.problem.Size(); i++ {
		result = append(result, i)
	}
	return result
}

func (p *Puzzle) ValueSet() []string {
	return p.valueSet
}

func (p *Puzzle) InvertSet() map[string]int {
	invertSet := map[string]int{}
	for i, v := range p.valueSet {
		invertSet[v] = i
	}
	return invertSet
}

func (p *Puzzle) Init(start string) {
	invertSet := p.InvertSet()
	for i, c := range start {
		if v, ok := invertSet[string(c)]; ok {
			p.problem.Set(i, v)
		}
	}
}

func (p *Puzzle) Solve(s csp.Settings) bool {
	return p.problem.Solve(s)
}

func (p *Puzzle) String() string {
	result := ""
	for i := 0; i < p.problem.Size(); i++ {
		v := " "
		if val := p.problem.Get(i).Value(); val >= 0 {
			v = p.valueSet[val]
		}
		result += v
	}
	return result
}
