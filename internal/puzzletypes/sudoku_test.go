package puzzletypes

import (
	"testing"

	"github.com/offpath/puzzleutils/internal/csp"
	"github.com/offpath/puzzleutils/internal/decide"
	"github.com/offpath/puzzleutils/internal/puzzle"
	"github.com/offpath/puzzleutils/internal/tracker"
)

var sudokuTests = []struct {
	name  string
	input string
	want  string
}{
	{
		name: "simple",
		input: "........." +
			".....3.85" +
			"..1.2...." +
			"...5.7..." +
			"..4...1.." +
			".9......." +
			"5......73" +
			"..2.1...." +
			"....4...9",
		want: "987654321\n" +
			"246173985\n" +
			"351928746\n" +
			"128537694\n" +
			"634892157\n" +
			"795461832\n" +
			"519286473\n" +
			"472319568\n" +
			"863745219",
	},
}

func TestSudoku(t *testing.T) {
	for _, tt := range sudokuTests {
		p := puzzle.NewPuzzle2()
		sudoku := NewSudoku(p)
		sudoku.Init(tt.input)
		settings := csp.Settings{
			DecisionTracker: tracker.PrintEveryLogN(10),
			Decider:         &decide.Min{},
		}
		if !p.Solve(settings) {
			t.Errorf("%s: failed to solve", tt.name)
		} else if sudoku.String() != tt.want {
			t.Errorf("%s: got %s, want %s", tt.name, sudoku.String(), tt.want)
		}
	}
}
