package puzzle

import (
	"testing"

	"github.com/offpath/puzzleutils/internal/csp"
	"github.com/offpath/puzzleutils/internal/decide"
)

// TODO(dneal): Sudoku unittests.

func TestNonogram(t *testing.T) {
	rows := [][]int{
		{8, 7, 5, 7},
		{5, 4, 3, 3},
		{3, 3, 2, 3},
		{4, 3, 2, 2},
		{3, 3, 2, 2},
		{3, 4, 2, 2},
		{4, 5, 2},
		{3, 5, 1},
		{4, 3, 2},
		{3, 4, 2},
		{4, 4, 2},
		{3, 6, 2},
		{3, 2, 3, 1},
		{4, 3, 4, 2},
		{3, 2, 3, 2},
		{6, 5},
		{4, 5},
		{3, 3},
		{3, 3},
		{1, 1},
	}
	cols := [][]int{
		{1},
		{1},
		{2},
		{4},
		{7},
		{9},
		{2, 8},
		{1, 8},
		{8},
		{1, 9},
		{2, 7},
		{3, 4},
		{6, 4},
		{8, 5},
		{1, 11},
		{1, 7},
		{8},
		{1, 4, 8},
		{6, 8},
		{4, 7},
		{2, 4},
		{1, 4},
		{5},
		{1, 4},
		{1, 5},
		{7},
		{5},
		{3},
		{1},
		{1},
	}
	nonogram := NewNonogramPuzzle(rows, cols)
	nonogram.problem.Solve(csp.Settings{Decider: &decide.First{}})
	got := nonogram.String()
	want := `XXXXXXXX.XXXXXXX.XXXXX.XXXXXXX
..XXXXX...XXXX....XXX....XXX..
...XXX.....XXX....XX.....XXX..
...XXXX.....XXX...XX.....XX...
....XXX.....XXX..XX......XX...
....XXX.....XXXX.XX.....XX....
....XXXX.....XXXXX......XX....
.....XXX.....XXXXX......X.....
.....XXXX.....XXX......XX.....
......XXX.....XXXX.....XX.....
......XXXX....XXXX....XX......
.......XXX...XXXXXX...XX......
.......XXX...XX.XXX...X.......
.......XXXX.XXX.XXXX.XX.......
........XXX.XX...XXX.XX.......
........XXXXXX...XXXXX........
.........XXXX....XXXXX........
.........XXX......XXX.........
.........XXX......XXX.........
..........X........X..........
`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s\n", got, want)
	}
}
