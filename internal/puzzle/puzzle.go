package puzzle

import (
	_ "fmt"
	
	"puzzleutils/internal/csp"
	"puzzleutils/internal/constraints"
)

type Puzzle struct {
	*csp.Problem
}

func NewPuzzle(size int, valueSet []string) *Puzzle {
	return &Puzzle{csp.NewProblem(size, valueSet)}
}

func (p *Puzzle) AllGroup() []int {
	var result []int
	for i := 0; i < p.Size(); i++ {
		result = append(result, i)
	}
	return result
}

func (p *Puzzle) InvertSet() map[string]int {
	invertSet := map[string]int{}
	for i, v := range p.ValueSet() {
		invertSet[v] = i
	}
	return invertSet
}

func (p *Puzzle) Init(start string) {
	invertSet := p.InvertSet()
	for i, c := range start {
		if v, ok := invertSet[string(c)]; ok {
			p.Set(i, v)
		}
	}
}

type GridEntry struct {
	Row int
	Col int
}

type GridPuzzle struct {
	*Puzzle
	width int
	height int
}

func NewGridPuzzle(width int, height int, valueSet []string) *GridPuzzle {
	return &GridPuzzle{NewPuzzle(width * height, valueSet), width, height}
}

func (p *GridPuzzle) AddGroup(group []GridEntry, constraint csp.ConstraintChecker) {
	var flatGroup []int
	for _, e := range group {
		flatGroup = append(flatGroup, e.Row * p.width + e.Col)
	}
	p.Puzzle.AddGroup(flatGroup, constraint)
}

func (p *GridPuzzle) RectGroup(row, col, height, width int) []GridEntry {
	var result []GridEntry
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			result = append(result, GridEntry{row+i,col+j})
		}
	}
	return result
}

func (p *GridPuzzle) ColumnGroups() [][]GridEntry {
	var result [][]GridEntry
	for i := 0; i < p.width; i++ {
		result = append(result, p.RectGroup(0, i, p.height, 1))
	}
	return result
}

func (p *GridPuzzle) RowGroups() [][]GridEntry {
	var result [][]GridEntry
	for i := 0; i < p.height; i++ {
		result = append(result, p.RectGroup(i, 0, 1, p.width))
	}
	return result
}

func NewSudokuPuzzle() *GridPuzzle {
	p := NewGridPuzzle(9, 9, []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"})
	for _, g := range p.ColumnGroups() {
		p.AddGroup(g, constraints.UniqueCovering())
	}
	for _, g := range p.RowGroups() {
		p.AddGroup(g, constraints.UniqueCovering())
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			p.AddGroup(p.RectGroup(i*3, j*3, 3, 3), constraints.UniqueCovering())
		}
	}
	return p
}

type nonogramConstraint struct{
	lengths []int
}
func (c nonogramConstraint) Init(all []*csp.Decision, size int){}
func (c nonogramConstraint) Apply(all, dirty []*csp.Decision) bool {
	b := csp.NewBuildupSet(len(all))
	var f func(index int)
	f = func(index int, start int) {
		sum := len(c.lengths) - index - 1
		for i := index; i < len(c.lengths); i++ {
			sum += c.lengths[i]
		}
		for i := start; i <= len(all) - sum; i++ {
			
		}
	}
	f(0)
	return true
}

func NewNonogramPuzzle(rows, cols [][]int) {
	p := NewGridPuzzle(len(cols), len(rows), []string{" ", "X"})
	for i, g := range p.RowGroups() {
		p.AddGroup(g, nonogramConstraint{rows[i]})
	}
	for i, g := range p.ColumnGroups() {
		p.AddGroup(g, nonogramConstraint{cols[i]})
	}
}
