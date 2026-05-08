package puzzle

import (
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
