package puzzletypes

import (
	"github.com/offpath/puzzleutils/internal/constraints2"
	"github.com/offpath/puzzleutils/internal/layouts"
	"github.com/offpath/puzzleutils/internal/puzzle"
)

func NewSudoku(p *puzzle.Puzzle2) *layouts.Grid {
	result := layouts.NewGrid(p, 9, 9)
	vs := p.GetIntRange(1, 9)
	result.Fill(vs)
	for i := 0; i < 9; i++ {
		p.AddConstraint(constraints2.NewUniqueConstraint(result.GetRow(i), vs))
		p.AddConstraint(constraints2.NewUniqueConstraint(result.GetCol(i), vs))
	}
	for i := range 3 {
		for j := range 3 {
			var block []*puzzle.Variable
			for k := range 3 {
				for l := range 3 {
					block = append(block, result.Get(i*3+k, j*3+l))
				}
			}
			p.AddConstraint(constraints2.NewUniqueConstraint(block, vs))
		}
	}
	return result
}
