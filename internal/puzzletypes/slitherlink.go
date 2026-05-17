package puzzletypes

import (
	"strings"

	"github.com/offpath/puzzleutils/internal/constraints2"
	"github.com/offpath/puzzleutils/internal/layouts"
	"github.com/offpath/puzzleutils/internal/puzzle"
)

type Slitherlink struct {
	graph *layouts.SquareGridGraph
	rows  int
	cols  int
}

func NewSlitherlink(p *puzzle.Puzzle2, input string) *Slitherlink {
	lines := strings.Split(input, "\n")
	rows := len(lines)
	cols := len(lines[0])

	g := layouts.NewSquareGridGraph(p, rows+1, cols+1)

	for r := 0; r <= rows; r++ {
		for c := 0; c <= cols; c++ {
			p.AddConstraint(constraints2.NewGridGraphPointConstraint(p, g.Get(r, c)))
		}
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if lines[r][c] != '.' {
				n := int(lines[r][c] - '0')
				p.AddConstraint(constraints2.NewGridGraphBoxConstraint(p, g.SurroundingArcs(r, c), n))
			}
		}
	}

	p.AddConstraint(constraints2.NewGridGraphLoopConstraint(p, g.Graph()))

	return &Slitherlink{
		graph: g,
		rows:  rows,
		cols:  cols,
	}
}

func (s *Slitherlink) horizontalString(row int) string {
	result := "."
	for col := 0; col < s.cols; col++ {
		arc := s.graph.Get(row, col).ArcTo(s.graph.Get(row, col+1))
		if arc != nil && arc.Var.Values().Value() != nil && arc.Var.Values().Value().Bool() {
			result += "-"
		} else {
			result += "X"
		}
		result += "."
	}
	return result
}

func (s *Slitherlink) verticalString(row int) string {
	result := ""
	for col := 0; col <= s.cols; col++ {
		arc := s.graph.Get(row, col).ArcTo(s.graph.Get(row+1, col))
		if arc != nil && arc.Var.Values().Value() != nil && arc.Var.Values().Value().Bool() {
			result += "|"
		} else {
			result += "X"
		}
		if col != s.cols {
			result += " "
		}
	}
	return result
}

func (s *Slitherlink) String() string {
	var result []string
	for row := 0; row < s.rows; row++ {
		result = append(result, s.horizontalString(row), s.verticalString(row))
	}
	result = append(result, s.horizontalString(s.rows))
	return strings.Join(result, "\n")
}
