package puzzle

import (
	"github.com/offpath/puzzleutils/internal/csp"
)

type Value struct {
	i int
}

func (val *Value) Int() int {
	return val.i
}

type ValueSet map[*Value]bool

func (vs ValueSet) Clone() ValueSet {
	result := map[*Value]bool{}
	for k, v := range vs {
		result[k] = v
	}
	return result
}

type Variable struct {
	p  *Puzzle2
	vs ValueSet
}

func (v *Variable) Set(vs ValueSet) {
	v.vs = vs
}

func (v *Variable) Value(val *Value) {
	v.vs = ValueSet{val: true}
}

type Puzzle2 struct {
	intValues map[int]*Value
	variables []*Variable
	// TODO(dneal): Groups
}

func NewPuzzle2() *Puzzle2 {
	return &Puzzle2{
		intValues: map[int]*Value{},
		variables: nil,
	}
}

func (p *Puzzle2) GetIntValue(i int) *Value {
	if val, ok := p.intValues[i]; ok {
		return val
	}
	val := &Value{i}
	p.intValues[i] = val
	return val
}

func (p *Puzzle2) GetIntRange(min, max int) ValueSet {
	result := ValueSet{}
	for i := min; i <= max; i++ {
		result[p.GetIntValue(i)] = true
	}
	return result
}

func (p *Puzzle2) NewVariable() *Variable {
	result := &Variable{
		p:  p,
		vs: ValueSet{},
	}
	p.variables = append(p.variables, result)
	return result
}

type Grid struct {
	rows      int
	cols      int
	variables [][]*Variable
}

func NewGrid(p *Puzzle2, rows int, cols int) *Grid {
	result := &Grid{
		rows:      rows,
		cols:      cols,
		variables: nil,
	}
	for i := 0; i < rows; i++ {
		var row []*Variable
		for j := 0; j < cols; j++ {
			row = append(row, p.NewVariable())
		}
		result.variables = append(result.variables, row)
	}
	return result
}

func (g *Grid) Fill(vs ValueSet) {
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			g.variables[i][j].Set(vs.Clone())
		}
	}
}

func (g *Grid) Get(row int, col int) *Variable {
	return g.variables[row][col]
}

func NewSudoku(p *Puzzle2) *Grid {
	result := NewGrid(p, 9, 9)
	result.Fill(p.GetIntRange(1, 9))
	// TODO(dneal): Groups
	return result
}

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

func (p *Puzzle) ToString() string {
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
