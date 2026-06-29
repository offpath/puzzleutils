package puzzletypes

import (
	"testing"

	"github.com/offpath/puzzleutils/internal/csp"
	"github.com/offpath/puzzleutils/internal/decide"
	"github.com/offpath/puzzleutils/internal/puzzle"
	"github.com/offpath/puzzleutils/internal/tracker"
)

var nonogramTests = []struct {
	name string
	rows [][]int
	cols [][]int
	want string
}{
	{
		name: "5x5 simple",
		rows: [][]int{
			{3},
			{1, 1},
			{5},
			{1, 1},
			{3},
		},
		cols: [][]int{
			{3},
			{1, 1},
			{5},
			{1, 1},
			{3},
		},
		want: "  XXX\n" +
			"  X X\n" +
			"XXXXX\n" +
			"X X  \n" +
			"XXX  ",
	},
	{
		name: "20x30 puzzle",
		rows: [][]int{
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
		},
		cols: [][]int{
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
		},
		want: "XXXXXXXX XXXXXXX XXXXX XXXXXXX\n" +
			"  XXXXX   XXXX    XXX    XXX  \n" +
			"   XXX     XXX    XX     XXX  \n" +
			"   XXXX     XXX   XX     XX   \n" +
			"    XXX     XXX  XX      XX   \n" +
			"    XXX     XXXX XX     XX    \n" +
			"    XXXX     XXXXX      XX    \n" +
			"     XXX     XXXXX      X     \n" +
			"     XXXX     XXX      XX     \n" +
			"      XXX     XXXX     XX     \n" +
			"      XXXX    XXXX    XX      \n" +
			"       XXX   XXXXXX   XX      \n" +
			"       XXX   XX XXX   X       \n" +
			"       XXXX XXX XXXX XX       \n" +
			"        XXX XX   XXX XX       \n" +
			"        XXXXXX   XXXXX        \n" +
			"         XXXX    XXXXX        \n" +
			"         XXX      XXX         \n" +
			"         XXX      XXX         \n" +
			"          X        X          ",
	},
}

func TestNonogram(t *testing.T) {
	for _, tt := range nonogramTests {
		p := puzzle.NewPuzzle()
		nonogram := NewNonogram(p, tt.rows, tt.cols)
		settings := csp.Settings{
			DecisionTracker: tracker.PrintEveryLogN(10),
			Decider:         &decide.First{},
		}
		if !p.Solve(settings) {
			t.Errorf("%s: failed to solve", tt.name)
		} else if nonogram.String() != tt.want {
			t.Errorf("%s: got \n%s\nwant \n%s", tt.name, nonogram.String(), tt.want)
		}
	}
}
