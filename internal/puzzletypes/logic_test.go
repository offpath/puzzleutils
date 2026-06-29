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
	{
		name: "complicated",
		input: `Chef:Freda,Karl,Sonia,Wade
Dish:cashew tofu,lemon snapper,smoked pork,turkey soup
Score:42,49,56,63

eq(Dish(Sonia), cashew tofu)
eq(Score(Freda), 49)
eq(Score(Karl), plus(Score(Chef(smoked pork)), 1))
eq(Score(Chef(turkey soup)), minus(Score(Sonia), 1))`,
		want: `Freda, smoked pork, 49
Karl, turkey soup, 56
Sonia, cashew tofu, 63
Wade, lemon snapper, 42`,
	},
}

func TestLogic(t *testing.T) {
	for _, tt := range logicTests {
		p := puzzle.NewPuzzle()
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
