package puzzletypes

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/offpath/puzzleutils/internal/csp"
	"github.com/offpath/puzzleutils/internal/decide"
	"github.com/offpath/puzzleutils/internal/puzzle"
	"github.com/offpath/puzzleutils/internal/trie"
)

var cryptogramTests = []struct {
	name  string
	input string
	want  string
}{
	{
		name:  "trivial",
		input: "ABBCCEEFEG",
		want:  "BOOKKEEPER",
	},
	{
		name:  "hard",
		input: "MTT OSQTGNGOSQUN QV PGD KQYU WSUJ SGJU MKU ZGZNUZNU RDW NGJU MKU XKUMWUK ZGZNUZNU WSMZ GWSUKN",
		want:  "ALL PHILOSOPHIES IF YOU RIDE THEM HOME ARE NONSENSE BUT SOME ARE GREATER NONSENSE THAN OTHERS",
	},
}

type cryptoTracker struct {
	solutions []string
	result    *CryptogramPuzzle
}

func (t *cryptoTracker) CaptureSolution(prob *csp.Problem) bool {
	var gotWords []string
	for _, group := range t.result.Words {
		var wordGot string
		for _, v := range group {
			if val := v.Values().Value(); val != nil {
				wordGot += val.Str()
			} else {
				wordGot += "?"
			}
		}
		gotWords = append(gotWords, wordGot)
	}
	t.solutions = append(t.solutions, strings.Join(gotWords, " "))
	return false
}

func TestCryptogram(t *testing.T) {
	tr := trie.New()
	tr.AddFile(filepath.Join("testdata", "ospd2.txt"))
	for _, tt := range cryptogramTests {
		p := puzzle.NewPuzzle()
		result := NewCryptogramPuzzle(p, tt.input, tr)
		tracker := &cryptoTracker{result: result}
		if !p.Solve(csp.Settings{Decider: &decide.First{}, SolutionTracker: tracker}) {
			t.Errorf("test: %s, failed to solve!\n", tt.name)
			continue
		}

		found := false
		for _, got := range tracker.solutions {
			if got == tt.want {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("test: %s, want: %s not found in solutions: %v\n", tt.name, tt.want, tracker.solutions)
		}
	}
}
