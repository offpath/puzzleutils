package puzzletypes

import (
	"testing"

	"github.com/offpath/puzzleutils/internal/csp"
	"github.com/offpath/puzzleutils/internal/decide"
	"github.com/offpath/puzzleutils/internal/puzzle"
	"github.com/offpath/puzzleutils/internal/tracker"
)

var logicTests = []struct {
	name  string
	input string
	want  string
}{
	{
		name: "basic",
		input: `Chef:Alice,Bob
Dish:Apple,Blueberry

eq(Chef(Apple), Bob)`,
		want: `Alice, Blueberry
Bob, Apple`,
	},
}

func TestLogic(t *testing.T) {
	for _, tt := range logicTests {
		p := puzzle.NewPuzzle2()
		logic := NewLogicPuzzle(p, tt.input)
		if !p.Solve(csp.Settings{Decider: &decide.First{}, DecisionTracker: tracker.PrintEveryN(1)}) {
			t.Errorf("test: %s, failed to solve!\n", tt.name)
		}
		got := logic.String()
		if got != tt.want {
			t.Errorf("test: %s, got: \n%s\nwant: \n%s\n", tt.name, got, tt.want)
		}
	}
}
