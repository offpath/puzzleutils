package puzzle

import (
	"fmt"
	"strings"

	"github.com/offpath/puzzleutils/internal/constraints"
	"github.com/offpath/puzzleutils/internal/csp"
	"github.com/offpath/puzzleutils/internal/trie"
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

var alphabet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func NewDropquotePuzzle(input string, t *trie.Trie) *Puzzle {
	lines := strings.Split(input, "\n")
	numCols := len(lines[0])
	numLetters := 0
	var cols [][]int
	var colSets []map[int]int
	var isCovering []bool
	for i := 0; i < numCols; i++ {
		cols = append(cols, nil)
	}
	var wordLengths []int
	currentLength := 0
	letterBanks := false
	for _, line := range lines {
		if line == "" && !letterBanks {
			letterBanks = true
			continue
		}
		if !letterBanks {
			for i, c := range line {
				if c == '.' {
					currentLength++
					cols[i] = append(cols[i], numLetters)
					numLetters++
				} else {
					if currentLength > 0 {
						wordLengths = append(wordLengths, currentLength)
						currentLength = 0
					}
				}
			}
		} else {
			s := map[int]int{}
			for _, c := range line {
				s[int(c)-int('A')]++
			}
			colSets = append(colSets, s)
			isCovering = append(isCovering, len(line) == len(cols[len(colSets)-1]))
			if len(colSets) == len(cols) {
				break
			}
		}
	}
	if currentLength > 0 {
		wordLengths = append(wordLengths, currentLength)
	}
	result := NewPuzzle(numLetters, alphabet)
	offset := 0
	for _, length := range wordLengths {
		var group []int
		for j := 0; j < length; j++ {
			group = append(group, offset+j)
		}
		result.problem.AddGroup(group, constraints.ValidWord(t, alphabet))
		offset += length
	}
	for i := 0; i < numCols; i++ {
		// TODO(dneal): Allow non-covering sets.
		result.problem.AddGroup(cols[i], constraints.SetCount(colSets[i], isCovering[i]))
	}
	return result
}
