package puzzle

import (
	"strings"

	"github.com/offpath/puzzleutils/internal/csp"
)

type slitherlinkPointConstraint struct{}

func (c slitherlinkPointConstraint) Init(all []*csp.Decision, size int) {}
func (c slitherlinkPointConstraint) Apply(all, dirty []*csp.Decision) bool {
	// Either 0 or 2 lines may connect to a point (no crossing lines)
	possible := 0
	set := 0
	for _, d := range all {
		if d.Value() == 1 {
			set++
		}
		if d.Possible(1) {
			possible++
		}
	}
	if set > 2 || (set == 1 && possible == 1) {
		return false
	}
	if set == 1 && possible == 2 {
		for _, d := range all {
			if d.Possible(1) {
				d.RestrictTo(1)
			}
		}
	} else if possible < 2 {
		for _, d := range all {
			d.RestrictTo(0)
		}
	}
	return true
}

type slitherlinkBoxConstraint struct{ n int }

func (c slitherlinkBoxConstraint) Init(all []*csp.Decision, size int) {
	// Optimize a bit by special-casing 0.
	if c.n == 0 {
		for _, d := range all {
			d.RestrictTo(0)
		}
	}
}
func (c slitherlinkBoxConstraint) Apply(all, dirty []*csp.Decision) bool {
	possible := 0
	set := 0
	for _, d := range all {
		if d.Value() == 1 {
			set++
		}
		if d.Possible(1) {
			possible++
		}
	}
	if possible < c.n || set > c.n {
		return false
	}
	if possible == c.n {
		for _, d := range all {
			if d.Possible(1) {
				d.RestrictTo(1)
			}
		}
		set = c.n
	}
	if set == c.n {
		for _, d := range all {
			if d.Value() != 1 {
				d.RestrictTo(0)
			}
		}
	}
	return true
}

type slitherlinkLoopConstraint struct{ g SlitherlinkPuzzle }

func (c slitherlinkLoopConstraint) Init(all []*csp.Decision, size int) {}
func (c slitherlinkLoopConstraint) Apply(all, dirty []*csp.Decision) bool {
	// Exactly one loop is allowed
	hasLoop := false
	hasUnfinishedLoop := false
	for i := 0; i < c.g.numRows; i++ {
		for j := 0; j < c.g.numCols; j++ {
			curPt := slitherlinkPoint{i, j}
			seenPts := map[slitherlinkPoint]bool{}
			prevLine := -1
			for {
				lineFound := -1
				for _, line := range c.g.pointToLines(curPt) {
					if line != prevLine && all[line].Value() == 1 {
						lineFound = line
						break
					}
				}
				if lineFound < 0 {
					if len(seenPts) > 0 {
						hasUnfinishedLoop = true
					}
					break
				}
				prevLine = lineFound
				seenPts[curPt] = true
				if curPt != c.g.lineToPointsMap[lineFound][0] {
					curPt = c.g.lineToPointsMap[lineFound][0]
				} else {
					curPt = c.g.lineToPointsMap[lineFound][1]
				}
				if seenPts[curPt] {
					hasLoop = true
					break
				}
			}
		}
	}
	return !(hasLoop && hasUnfinishedLoop)
}

type slitherlinkPoint struct {
	row, col int
}

type SlitherlinkPuzzle struct {
	*Puzzle
	numRows, numCols int
	lineToPointsMap  map[int][]slitherlinkPoint
}

func (g SlitherlinkPuzzle) numLines() int {
	return g.numCols*(g.numRows+1) + g.numRows*(g.numCols+1)
}

func (g SlitherlinkPuzzle) pointToLines(pt slitherlinkPoint) []int {
	var result []int
	if pt.row > 0 {
		result = append(result, g.verticalLine(slitherlinkPoint{pt.row - 1, pt.col}))
	}
	if pt.col > 0 {
		result = append(result, g.horizontalLine(slitherlinkPoint{pt.row, pt.col - 1}))
	}
	if pt.row < g.numRows {
		result = append(result, g.verticalLine(pt))
	}
	if pt.col < g.numCols {
		result = append(result, g.horizontalLine(pt))
	}
	return result
}

func (g SlitherlinkPuzzle) boxToLines(pt slitherlinkPoint) []int {
	return []int{g.horizontalLine(pt), g.verticalLine(pt), g.horizontalLine(slitherlinkPoint{pt.row + 1, pt.col}), g.verticalLine(slitherlinkPoint{pt.row, pt.col + 1})}
}

func (g SlitherlinkPuzzle) horizontalLine(pt slitherlinkPoint) int {
	// Horizontal lines come first.
	return pt.row*(g.numCols) + pt.col
}

func (g SlitherlinkPuzzle) verticalLine(pt slitherlinkPoint) int {
	// Vertical lines come after horizontal lines
	return g.numCols*(g.numRows+1) + pt.row*(g.numCols+1) + pt.col

}

func NewSlitherlinkPuzzle(input string) SlitherlinkPuzzle {
	lines := strings.Split(input, "\n")
	g := SlitherlinkPuzzle{
		numRows:         len(lines),
		numCols:         len(lines[0]),
		lineToPointsMap: map[int][]slitherlinkPoint{},
	}
	for i := 0; i <= g.numRows; i++ {
		for j := 0; j <= g.numCols; j++ {
			pt := slitherlinkPoint{i, j}
			for _, line := range g.pointToLines(pt) {
				g.lineToPointsMap[line] = append(g.lineToPointsMap[line], pt)
			}
		}
	}
	g.Puzzle = NewPuzzle(g.numLines(), []string{"0", "1"})
	for i := 0; i <= g.numRows; i++ {
		for j := 0; j <= g.numCols; j++ {
			g.problem.AddGroup(g.pointToLines(slitherlinkPoint{i, j}), slitherlinkPointConstraint{})
		}
	}
	for i := 0; i < g.numRows; i++ {
		for j := 0; j < g.numCols; j++ {
			if lines[i][j] != '.' {
				g.problem.AddGroup(g.boxToLines(slitherlinkPoint{i, j}), slitherlinkBoxConstraint{int(lines[i][j] - '0')})
			}
		}
	}
	var group []int
	for i := 0; i < g.numLines(); i++ {
		group = append(group, i)
	}
	g.problem.AddGroup(group, slitherlinkLoopConstraint{g})

	return g
}

func (g SlitherlinkPuzzle) horizontalString(pt slitherlinkPoint) string {
	result := "."
	for col := 0; col < g.numCols; col++ {
		if g.problem.Get(g.horizontalLine(slitherlinkPoint{pt.row, col})).Value() == 1 {
			result += "-"
		} else {
			result += "X"
		}
		result += "."
	}
	return result
}

func (g SlitherlinkPuzzle) verticalString(pt slitherlinkPoint) string {
	result := ""
	for col := 0; col <= g.numCols; col++ {
		if g.problem.Get(g.verticalLine(slitherlinkPoint{pt.row, col})).Value() == 1 {
			result += "|"
		} else {
			result += "X"
		}
		if col != g.numCols {
			result += " "
		}
	}
	return result
}

func (g SlitherlinkPuzzle) String() string {
	result := []string{}
	for row := 0; row < g.numRows; row++ {
		result = append(result, g.horizontalString(slitherlinkPoint{row, 0}), g.verticalString(slitherlinkPoint{row, 0}))
	}
	result = append(result, g.horizontalString(slitherlinkPoint{g.numRows, 0}))
	return strings.Join(result, "\n")
}
