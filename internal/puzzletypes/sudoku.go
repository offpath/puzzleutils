package puzzletypes

import (
	"github.com/offpath/puzzleutils/internal/constraints"
	"github.com/offpath/puzzleutils/internal/layouts"
	"github.com/offpath/puzzleutils/internal/puzzle"
)

func NewSudoku(p *puzzle.Puzzle) *layouts.Grid {
	result := layouts.NewGrid(p, 9, 9)
	vs := p.GetIntRange(1, 9)
	result.Fill(vs)
	for i := range 9 {
		p.AddConstraint(constraints.NewUniqueConstraint(result.GetRow(i), vs))
		p.AddConstraint(constraints.NewUniqueConstraint(result.GetCol(i), vs))
	}
	for i := range 3 {
		for j := range 3 {
			p.AddConstraint(constraints.NewUniqueConstraint(result.GetRect(i*3, j*3, 3, 3), vs))
		}
	}
	return result
}
