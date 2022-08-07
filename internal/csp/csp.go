// Package csp implements a recursive constraint satisfaction problem
// solver.
package csp

import (
	"fmt"
)

// Decision represents a single decision to be made, for example a
// single grid square in a sudoku or similar puzzle.
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

// Value returns the value of a decision if the decision has been made
// or -1 if multiple possibilities remain.
func (d *Decision) Value() int {
	if len(d.possibilities) != 1 {
		return -1
	}
	for k := range d.possibilities {
		return k
	}
	return -1
}

func (d *Decision) StringValue() string {
	if v := d.Value(); v < 0 {
		return ""
	} else {
		return d.p.valueSet[v]
	}
}

// Count returns the number of remaining options for the decision.
func (d *Decision) Count() int {
	return len(d.possibilities)
}

func (d *Decision) undo(item undoRestricts) {
	for _, v := range item {
		d.possibilities[v] = true
	}
}

// Restrict removes i as a possibility for the decision.
func (d *Decision) Restrict(i int) {
	if d.possibilities[i] {
		delete(d.possibilities, i)
		d.p.setDirty(d, i)
		if len(d.possibilities) == 0 {
			d.p.setConflict()
		}
	}
}

// RestrictTo removes everything but i as a possibility for the
// decision.
func (d *Decision) RestrictTo(i int) {
	for k := range d.possibilities {
		if k != i {
			d.Restrict(k)
		}
	}

}

func (d *Decision) RestrictToSet(s map[int]bool) {
	for k := range d.possibilities {
		if !s[k] {
			d.Restrict(k)
		}
	}
}

func (d *Decision) Possible(i int) bool {
	return d.possibilities[i]
}

// Group represents a grouping of Decisions over which to apply a
// constraint, for example a row or column in a sudoku puzzle with a
// uniqueness constraint.
type Group struct {
	decisions []*Decision
	constraint ConstraintChecker
}

func (g *Group) Decisions() []*Decision {
	return g.decisions
}

// A ConstraintChecker is any object that can be used to validate a
// constraint over a group.
type ConstraintChecker interface {
	Init(all []*Decision, size int)
	Apply(all, dirty []*Decision) bool
}

// A Printer is an object that will be called at certain critical
// points during a solve, for example when a solution has been
// found. In spite of the name, it may also simply memorize,
// summarize, sample, or do anything else really.
type Printer interface {
	MakeDecision(p *Problem)
	CaptureSolution(p *Problem)
}

// A Decider decides which decision to decide next. Seriously. Wow
// this needs a better name.
type Decider interface {
	Decide(d []*Decision, g []*Group) *Decision
}

type Settings interface {
	Printer
	Decider
}

type undoRestricts []int 

// A Problem captures the decisions and groups, as well as any
// ephemeral state used to solve the problem.
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
	if p.conflict {
		p.dirty = map[*Decision]bool{}
		p.conflict = false
		return false
	}
	for ; len(p.dirty) > 0; {
		groups := map[*Group][]*Decision{}
		for d, _ := range p.dirty {
			for _, g := range d.groups {
				groups[g] = append(groups[g], d)
			}
		}
		p.dirty = map[*Decision]bool{}
		for g, dirty := range groups {
			if !g.constraint.Apply(g.decisions, dirty) {
				p.conflict = true
				break
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

// Attempts to solve the Problem, and returns true if a solution
// exists.
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
	d := s.Decide(ds, p.groups)
	for i := 0; i < len(p.valueSet); i++ {  // Improve by iterating over available values for d?
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

// Print the current state of the problem. 
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

// Set the value for a given decision. Used to set the initial
// configuration, for example the givens on a sudoku problem.
func (p *Problem) Set(i int, val int) {
	p.decisions[i].RestrictTo(val)
}

func (p *Problem) AddGroup(group []int, constraint ConstraintChecker) {
	g := Group{constraint: constraint}
	for _, d := range group {
		g.decisions = append(g.decisions, p.decisions[d])
		p.decisions[d].groups = append(p.decisions[d].groups, &g)
	}
	constraint.Init(g.decisions, len(p.valueSet))
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

