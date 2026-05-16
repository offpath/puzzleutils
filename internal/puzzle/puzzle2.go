package puzzle

import (
	"github.com/offpath/puzzleutils/internal/csp"
)

type Value struct {
	val interface{}
	id  int
}

func (val *Value) Int() int {
	return val.val.(int)
}

func (val *Value) Bool() bool {
	return val.val.(bool)
}

func (val *Value) Str() string {
	return val.val.(string)
}

func (val *Value) Raw() interface{} {
	return val.val
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
	// rep is the parent node in the disjoint-set of equal variables.
	// Variables that are marked equal will share the same root representative.
	rep *Variable
}

func (v *Variable) root() *Variable {
	root := v
	for root.rep != root {
		root = root.rep
	}
	return root
}

// MarkEqual marks the variable as equal to all of the provided variables.
func (v *Variable) MarkEqual(vars ...*Variable) {
	v.p.MarkEqual(append([]*Variable{v}, vars...)...)
}

func (v *Variable) SetValueRange(vs ValueSet) {
	v = v.root()
	v.possibleValues = vs
	v.currentValues = vs.Clone()
}

func (v *Variable) SetInitialValue(val *Value) {
	v = v.root()
	v.possibleValues = ValueSet{val: true}
	v.currentValues = ValueSet{val: true}
}

func (v *Variable) EliminateValue(val *Value) {
	v = v.root()
	delete(v.currentValues, val)
	v.p.problem.Get(v.id).Restrict(val.id)
}

func (v *Variable) RestrictTo(val *Value) {
	v = v.root()
	v.currentValues = ValueSet{val: true}
	v.p.problem.Get(v.id).RestrictTo(val.id)
}

func (v *Variable) Values() ValueSet {
	return v.root().currentValues
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
	values    map[interface{}]*Value
	variables []*Variable
	groups    []*group
	problem   *csp.Problem
}

func NewPuzzle2() *Puzzle2 {
	return &Puzzle2{
		values:    map[interface{}]*Value{},
		variables: nil,
		groups:    nil,
		problem:   nil,
	}
}

func (p *Puzzle2) GetIntValue(i int) *Value {
	if val, ok := p.values[i]; ok {
		return val
	}
	val := &Value{val: i, id: len(p.values)}
	p.values[i] = val
	return val
}

func (p *Puzzle2) GetBoolValue(b bool) *Value {
	if val, ok := p.values[b]; ok {
		return val
	}
	val := &Value{val: b, id: len(p.values)}
	p.values[b] = val
	return val
}

func (p *Puzzle2) GetStringValue(s string) *Value {
	if val, ok := p.values[s]; ok {
		return val
	}
	val := &Value{val: s, id: len(p.values)}
	p.values[s] = val
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
	result.rep = result
	p.variables = append(p.variables, result)
	return result
}

func (p *Puzzle2) AddConstraint(c Constraint) {
	p.groups = append(p.groups, &group{
		variables: c.Variables(),
		c:         c,
	})
}

func (p *Puzzle2) MarkEqual(vars ...*Variable) {
	if len(vars) < 2 {
		return
	}
	first := vars[0].root()
	for _, v := range vars[1:] {
		rep := v.root()
		if rep != first {
			rep.rep = first

			if len(first.possibleValues) == 0 {
				first.possibleValues = rep.possibleValues.Clone()
			} else if len(rep.possibleValues) > 0 {
				for val := range first.possibleValues {
					if !rep.possibleValues[val] {
						delete(first.possibleValues, val)
					}
				}
			}
			first.currentValues = first.possibleValues.Clone()
		}
	}
}

type constraintShim struct {
	c Constraint
}

func (c constraintShim) Init(all []*csp.Decision, size int) {}
func (c constraintShim) Apply(all, dirty []*csp.Decision) bool {
	for _, v := range c.c.Variables() {
		v = v.root()
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
	numDecisions := 0
	for _, v := range p.variables {
		if v.rep == v {
			v.id = numDecisions
			numDecisions++
		}
	}
	for _, v := range p.variables {
		v.id = v.root().id
	}

	p.problem = csp.NewProblem(numDecisions, len(p.values))
	for _, v := range p.variables {
		if v.rep == v {
			group := []int{v.id}
			p.problem.AddGroup(group, valueSetConstraint{v.possibleValues})
		}
	}
	for _, g := range p.groups {
		var group []int
		seen := map[int]bool{}
		for _, v := range g.variables {
			if !seen[v.id] {
				group = append(group, v.id)
				seen[v.id] = true
			}
		}
		p.problem.AddGroup(group, constraintShim{g.c})
	}
	return p.problem.Solve(settings)
}
