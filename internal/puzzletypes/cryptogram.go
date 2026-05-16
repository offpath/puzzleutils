package puzzletypes

import (
	"strings"

	"github.com/offpath/puzzleutils/internal/constraints2"
	"github.com/offpath/puzzleutils/internal/puzzle"
	"github.com/offpath/puzzleutils/internal/trie"
)

var alphabet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func NewCryptogramPuzzle(p *puzzle.Puzzle2, input string, t *trie.Trie) [][]*puzzle.Variable {
	words := strings.Split(input, " ")
	uniqueLetters := map[string]*puzzle.Variable{}

	vs := puzzle.ValueSet{}
	for _, letter := range alphabet {
		vs[p.GetStringValue(letter)] = true
	}

	for _, word := range words {
		for _, c := range word {
			charStr := string(c)
			if _, ok := uniqueLetters[charStr]; !ok {
				v := p.NewVariable()
				v.SetValueRange(vs)
				uniqueLetters[charStr] = v
			}
		}
	}

	var allVars []*puzzle.Variable
	for _, v := range uniqueLetters {
		allVars = append(allVars, v)
	}
	p.AddConstraint(constraints2.NewUniqueConstraint(allVars, vs))

	var result [][]*puzzle.Variable
	for _, word := range words {
		var group []*puzzle.Variable
		for _, c := range word {
			group = append(group, uniqueLetters[string(c)])
		}
		p.AddConstraint(constraints2.NewValidWordConstraint(p, group, t, alphabet))
		result = append(result, group)
	}

	return result
}
