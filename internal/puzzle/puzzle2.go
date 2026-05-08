package puzzle

import (
	"strconv"

	"github.com/offpath/puzzleutils/internal/csp"
)

type Value struct {
	i  int
	id int
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

func (vs ValueSet) Value() *Value {
	if len(vs) == 1 {
		for k := range vs {
			return k
		}
	}
	return nil
}

type Variable struct {
	p              *Puzzle2
	possibleValues ValueSet
	currentValues  ValueSet
	id             int
}

func (v *Variable) SetValueRange(vs ValueSet) {
	v.possibleValues = vs
	v.currentValues = vs.Clone()
}

func (v *Variable) SetInitialValue(val *Value) {
	v.possibleValues = ValueSet{val: true}
	v.currentValues = ValueSet{val: true}
}

func (v *Variable) EliminateValue(val *Value) {
	delete(v.currentValues, val)
	v.p.problem.Get(v.id).Restrict(val.id)
}

func (v *Variable) RestrictTo(val *Value) {
	v.currentValues = ValueSet{val: true}
	v.p.problem.Get(v.id).RestrictTo(val.id)
}

func (v *Variable) Values() ValueSet {
	return v.currentValues
}

type Constraint interface {
	Variables() []*Variable
	Check() bool
}

type group struct {
	variables []*Variable
	c         Constraint
}

type Puzzle2 struct {
	intValues map[int]*Value
	variables []*Variable
	groups    []*group
	problem   *csp.Problem
}

func NewPuzzle2() *Puzzle2 {
	return &Puzzle2{
		intValues: map[int]*Value{},
		variables: nil,
		groups:    nil,
		problem:   nil,
	}
}

func (p *Puzzle2) GetIntValue(i int) *Value {
	if val, ok := p.intValues[i]; ok {
		return val
	}
	val := &Value{i, len(p.intValues)}
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
		p:              p,
		possibleValues: ValueSet{},
		currentValues:  ValueSet{},
		id:             len(p.variables),
	}
	p.variables = append(p.variables, result)
	return result
}

func (p *Puzzle2) AddConstraint(c Constraint) {
	p.groups = append(p.groups, &group{
		variables: c.Variables(),
		c:         c,
	})
}

type constraintShim struct {
	c Constraint
}

func (c constraintShim) Init(all []*csp.Decision, size int) {}
func (c constraintShim) Apply(all, dirty []*csp.Decision) bool {
	for _, v := range c.c.Variables() {
		v.currentValues = v.possibleValues.Clone()
		for val := range v.currentValues {
			if !v.p.problem.Get(v.id).Possible(val.id) {
				delete(v.currentValues, val)
			}
		}
	}
	return c.c.Check()
}

type valueSetConstraint struct {
	vs ValueSet
}

func (c valueSetConstraint) Init(all []*csp.Decision, size int) {
	for _, d := range all {
		set := map[int]bool{}
		for v := range c.vs {
			set[v.id] = true
		}
		d.RestrictToSet(set)
	}
}

func (c valueSetConstraint) Apply(all, dirty []*csp.Decision) bool {
	return true
}

func (p *Puzzle2) Solve(settings csp.Settings) bool {
	p.problem = csp.NewProblem(len(p.variables), len(p.intValues))
	for _, v := range p.variables {
		group := []int{v.id}
		p.problem.AddGroup(group, valueSetConstraint{v.possibleValues})
	}
	for _, g := range p.groups {
		var group []int
		for _, v := range g.variables {
			group = append(group, v.id)
		}
		p.problem.AddGroup(group, constraintShim{g.c})
	}
	return p.problem.Solve(settings)
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
			g.variables[i][j].SetValueRange(vs.Clone())
		}
	}
}

func (g *Grid) Get(row int, col int) *Variable {
	return g.variables[row][col]
}

func (g *Grid) GetRow(row int) []*Variable {
	return g.variables[row]
}

func (g *Grid) GetCol(col int) []*Variable {
	var result []*Variable
	for i := 0; i < g.rows; i++ {
		result = append(result, g.variables[i][col])
	}
	return result
}

type UniqueConstraint struct {
	variables []*Variable
	vs        ValueSet
}

func (c *UniqueConstraint) Variables() []*Variable { return c.variables }
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
	if missingCount > len(c.vs)-len(c.variables) {
		return false
	}
	return true
}

func NewSudoku(p *Puzzle2) *Grid {
	result := NewGrid(p, 9, 9)
	vs := p.GetIntRange(1, 9)
	result.Fill(vs)
	for i := 0; i < 9; i++ {
		p.AddConstraint(&UniqueConstraint{result.GetRow(i), vs})
		p.AddConstraint(&UniqueConstraint{result.GetCol(i), vs})
	}
	for i := range 3 {
		for j := range 3 {
			var block []*Variable
			for k := range 3 {
				for l := range 3 {
					block = append(block, result.Get(i*3+k, j*3+l))
				}
			}
			p.AddConstraint(&UniqueConstraint{block, vs})
		}
	}
	return result
}

func (g *Grid) Init(start string) {
	for i, r := range g.variables {
		for j, v := range r {
			if start[i*g.cols+j] != '.' {
				v.SetInitialValue(v.p.GetIntValue(int(start[i*g.cols+j] - '0')))
			}
		}
	}
}

func (g *Grid) String() string {
	result := ""
	for i, r := range g.variables {
		for _, v := range r {
			if val := v.Values().Value(); val != nil {
				result += strconv.Itoa(val.Int())
			} else {
				result += "."
			}
		}
		if i < g.rows-1 {
			result += "\n"
		}
	}
	return result
}
