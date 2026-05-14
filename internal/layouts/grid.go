package layouts

import (
	"strconv"

	"github.com/offpath/puzzleutils/internal/puzzle"
)

type Grid struct {
	rows      int
	cols      int
	variables [][]*puzzle.Variable
	p         *puzzle.Puzzle2
}

func NewGrid(p *puzzle.Puzzle2, rows int, cols int) *Grid {
	result := &Grid{
		rows:      rows,
		cols:      cols,
		variables: nil,
		p:         p,
	}
	for i := 0; i < rows; i++ {
		var row []*puzzle.Variable
		for j := 0; j < cols; j++ {
			row = append(row, p.NewVariable())
		}
		result.variables = append(result.variables, row)
	}
	return result
}

func (g *Grid) Fill(vs puzzle.ValueSet) {
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			g.variables[i][j].SetValueRange(vs.Clone())
		}
	}
}

func (g *Grid) Get(row int, col int) *puzzle.Variable {
	return g.variables[row][col]
}

func (g *Grid) GetRow(row int) []*puzzle.Variable {
	return g.variables[row]
}

func (g *Grid) GetCol(col int) []*puzzle.Variable {
	var result []*puzzle.Variable
	for i := 0; i < g.rows; i++ {
		result = append(result, g.variables[i][col])
	}
	return result
}

func (g *Grid) GetRect(row int, col int, height int, width int) []*puzzle.Variable {
	var result []*puzzle.Variable
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			result = append(result, g.variables[row+i][col+j])
		}
	}
	return result
}

func (g *Grid) Init(start string) {
	for i, r := range g.variables {
		for j, v := range r {
			if start[i*g.cols+j] != '.' {
				v.SetInitialValue(g.p.GetIntValue(int(start[i*g.cols+j] - '0')))
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
