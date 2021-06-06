package csp

import (
	"fmt"
)

type Decision struct {
	possibilities map[int]bool
	groups []*Group
	p *Problem
}

func newDecision(size int, p *Problem) *Decision {
	result := Decision{
		possibilities: map[int]bool{},
		p: p,
	}
	for i := 0; i < size; i++ {
		result.possibilities[i] = true
	}
	return &result
}

func (d *Decision) Value() int {
	if len(d.possibilities) != 1 {
		return -1
	}
	for k := range d.possibilities {
		return k
	}
	return -1
}

func (d *Decision) Count() int {
	return len(d.possibilities)
}

func (d *Decision) undo(item undoRestricts) {
	for _, v := range item {
		d.possibilities[v] = true
	}
}

func (d *Decision) Restrict(i int) {
	if d.possibilities[i] {
		delete(d.possibilities, i)
		d.p.setDirty(d, i)
		if len(d.possibilities) == 0 {
			d.p.setConflict()
		}
	}
}

func (d *Decision) RestrictTo(i int) {
	for k := range d.possibilities {
		if k != i {
			d.Restrict(k)
		}
	}

}

type Group struct {
	decisions []*Decision
	constraint ConstraintChecker
}

type ConstraintChecker interface {
	Check(all, dirty []*Decision) bool
}

type Printer interface {
	MakeDecision(p *Problem)
	CaptureSolution(p *Problem)
}

type Decider interface {
	Decide(d []*Decision) *Decision
}

type Settings interface {
	Printer
	Decider
}

type undoRestricts []int 

type Problem struct {
	valueSet []string
	decisions []*Decision
	groups []*Group
	undoStack []map[*Decision]undoRestricts
	dirty map[*Decision]bool
	conflict bool
}

func (p *Problem) Size() int {
	return len(p.decisions)
}

func (p *Problem) ValueSet() []string {
	return p.valueSet
}

func (p *Problem) check() bool {
	//fmt.Printf("Check\n")
	if p.conflict {
		p.dirty = map[*Decision]bool{}
		p.conflict = false
		return false
	}
	for ; len(p.dirty) > 0; {
		//fmt.Printf("Len = %d\n", len(p.dirty))
		groups := map[*Group][]*Decision{}
		for d, _ := range p.dirty {
			for _, g := range d.groups {
				groups[g] = append(groups[g], d)
			}
		}
		p.dirty = map[*Decision]bool{}
		for g, dirty := range groups {
			if !g.constraint.Check(g.decisions, dirty) {
				p.dirty = map[*Decision]bool{}
				p.conflict = false
				return false
			}
		}
		if p.conflict {
			p.dirty = map[*Decision]bool{}
			p.conflict = false
			return false
		}
	}
	p.conflict = false
	return true
}

func (p *Problem) snapshot() {
	p.undoStack = append(p.undoStack, map[*Decision]undoRestricts{})
}

func (p *Problem) undo() {
	i := len(p.undoStack) - 1
	m := p.undoStack[i]
	p.undoStack = p.undoStack[:i]
	for d, item := range m {
		d.undo(item)
	}
}

func (p *Problem) setDirty(d *Decision, restrict int) {
	if i := len(p.undoStack) - 1; i >= 0 {
		p.undoStack[i][d] = append(p.undoStack[i][d], restrict)
	}
	p.dirty[d] = true
}

func (p *Problem) setConflict() {
	p.conflict = true
}

func (p *Problem) Solve(s Settings) bool {
	if !p.check() {
		return false
	}
	return p.recSolve(s)
}

func (p *Problem) recSolve(s Settings) bool {
	var ds []*Decision
	for _, d := range p.decisions {
		if d.Value() == -1 {
			ds = append(ds, d)
		}
	}
	if len(ds) == 0 {
		s.CaptureSolution(p)
		return true
	}
	d := s.Decide(ds)
	/*d := 0
	for ; d < len(p.decisions) && p.decisions[d].Value() >= 0; d++ {}
	if d == len(p.decisions) {
		s.CaptureSolution(p)
		return true
	}*/
	for i := 0; i < len(p.valueSet); i++ {
		//fmt.Printf("Trying (%d, %d)\n", d, i)
		s.MakeDecision(p)
		p.snapshot()
		d.RestrictTo(i)
		if p.check() && p.recSolve(s) {
			return true
		}
		p.undo()
	}
	return false
}

func (p *Problem) Print() {
	for i, d := range p.decisions {
		val := ""
		if v := d.Value(); v >= 0 {
			val = p.valueSet[v]
		}
		fmt.Printf("%d: %s\n", i, val)
	}
	fmt.Println("")
}

func (p *Problem) Set(i int, val int) {
	p.decisions[i].RestrictTo(val)
}

func (p *Problem) AddGroup(group []int, constraint ConstraintChecker) {
	g := Group{constraint: constraint}
	for _, d := range group {
		g.decisions = append(g.decisions, p.decisions[d])
		p.decisions[d].groups = append(p.decisions[d].groups, &g)
	}
	p.groups = append(p.groups, &g)
}

func NewProblem(size int, valueSet []string) *Problem {
	p := Problem{
		valueSet: valueSet,
		dirty: map[*Decision]bool{},
	}
	for i := 0; i < size; i++ {
		p.decisions = append(p.decisions, newDecision(len(valueSet), &p))
	}
	return &p
}

