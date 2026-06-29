package puzzletypes

import (
	"github.com/offpath/puzzleutils/internal/constraints2"
	"github.com/offpath/puzzleutils/internal/layouts"
	"github.com/offpath/puzzleutils/internal/puzzle"
)

func NewNonogram(p *puzzle.Puzzle, rows, cols [][]int) *layouts.Grid {
	result := layouts.NewGrid(p, len(rows), len(cols))
	vs := puzzle.ValueSet{
		p.GetBoolValue(false): true,
		p.GetBoolValue(true):  true,
	}
	result.Fill(vs)
	for i, lengths := range rows {
		p.AddConstraint(constraints2.NewNonogramConstraint(p, result.GetRow(i), lengths))
	}
	for i, lengths := range cols {
		p.AddConstraint(constraints2.NewNonogramConstraint(p, result.GetCol(i), lengths))
	}
	return result
}
