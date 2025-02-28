package puzzle

import (
	"fmt"

	"github.com/offpath/puzzleutils/internal/constraints"
	"github.com/offpath/puzzleutils/internal/csp"
)

type GridEntry struct {
	Row int
	Col int
}

type GridPuzzle struct {
	*Puzzle
	width  int
	height int
}

func NewGridPuzzle(width int, height int, valueSet []string) *GridPuzzle {
	return &GridPuzzle{NewPuzzle(width*height, valueSet), width, height}
}

func (p *GridPuzzle) AddGroup(group []GridEntry, constraint csp.ConstraintChecker) {
	var flatGroup []int
	for _, e := range group {
		flatGroup = append(flatGroup, e.Row*p.width+e.Col)
	}
	p.Puzzle.problem.AddGroup(flatGroup, constraint)
}

func (p *GridPuzzle) RectGroup(row, col, height, width int) []GridEntry {
	var result []GridEntry
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			result = append(result, GridEntry{row + i, col + j})
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

func (p *GridPuzzle) ToString() string {
	result := ""
	for i := 0; i < p.height; i++ {
		for j := 0; j < p.width; j++ {
			v := " "
			if val := p.problem.Get(i*p.width + j).Value(); val >= 0 {
				v = p.valueSet[val]
			}
			result += v
		}
		result += "\n"
	}
	return result
}

func (p *GridPuzzle) Print() {
	fmt.Print(p.ToString())
}

func NewSudokuPuzzle() *GridPuzzle {
	p := NewGridPuzzle(9, 9, []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"})
	for _, g := range p.ColumnGroups() {
		p.AddGroup(g, constraints.Unique(true))
	}
	for _, g := range p.RowGroups() {
		p.AddGroup(g, constraints.Unique(true))
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			p.AddGroup(p.RectGroup(i*3, j*3, 3, 3), constraints.Unique(true))
		}
	}
	return p
}

type nonogramConstraint struct {
	lengths []int
}

func (c nonogramConstraint) Init(all []*csp.Decision, size int) {}
func (c nonogramConstraint) Apply(all, dirty []*csp.Decision) bool {
	b := constraints.NewBuildupSet(len(all))
	var f func(lengths []int, ds []*csp.Decision)
	f = func(lengths []int, ds []*csp.Decision) {
		if len(ds) == 0 {
			return
		}
		sum := len(lengths) - 1
		for _, n := range lengths {
			sum += n
		}
		if sum > len(ds) {
			return
		}
		if ds[0].Possible(0) {
			b.Push(0)
			f(lengths, ds[1:])
			b.Pop()
		}
		if len(lengths) == 0 {
			return
		}
		for i := 0; i < lengths[0]; i++ {
			if !ds[i].Possible(1) {
				return
			}
			b.Push(1)
			defer b.Pop()
		}
		if len(ds) > lengths[0] {
			if !ds[lengths[0]].Possible(0) {
				return
			}
			b.Push(0)
			defer b.Pop()
			f(lengths[1:], ds[lengths[0]+1:])
		}
	}
	f(c.lengths, all)
	b.Export(all)
	return true
}

func NewNonogramPuzzle(rows, cols [][]int) *GridPuzzle {
	p := NewGridPuzzle(len(cols), len(rows), []string{".", "X"})
	for i, g := range p.RowGroups() {
		p.AddGroup(g, nonogramConstraint{rows[i]})
	}
	for i, g := range p.ColumnGroups() {
		p.AddGroup(g, nonogramConstraint{cols[i]})
	}
	return p
}
