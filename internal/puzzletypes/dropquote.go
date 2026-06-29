package puzzletypes

import (
	"strings"

	"github.com/offpath/puzzleutils/internal/constraints"
	"github.com/offpath/puzzleutils/internal/puzzle"
	"github.com/offpath/puzzleutils/internal/trie"
)

type Dropquote struct {
	variables []*puzzle.Variable
}

func NewDropquote(p *puzzle.Puzzle, input string, t *trie.Trie) *Dropquote {
	lines := strings.Split(input, "\n")
	numCols := len(lines[0])

	var cols [][]*puzzle.Variable
	var colSets []map[*puzzle.Value]int
	var isCovering []bool

	for i := 0; i < numCols; i++ {
		cols = append(cols, nil)
	}

	var allVars []*puzzle.Variable
	var wordGroups [][]*puzzle.Variable
	var currentWord []*puzzle.Variable
	letterBanks := false

	alphabet := p.GetAlphabet()

	for _, line := range lines {
		if line == "" && !letterBanks {
			letterBanks = true
			continue
		}
		if !letterBanks {
			for i, c := range line {
				if c == '.' {
					v := p.NewVariable()
					v.SetValueRange(alphabet)
					currentWord = append(currentWord, v)
					cols[i] = append(cols[i], v)
					allVars = append(allVars, v)
				} else {
					if len(currentWord) > 0 {
						wordGroups = append(wordGroups, currentWord)
						currentWord = nil
					}
				}
			}
		} else {
			s := map[*puzzle.Value]int{}
			for _, c := range line {
				s[p.GetStringValue(string(c))]++
			}
			colSets = append(colSets, s)
			isCovering = append(isCovering, len(line) == len(cols[len(colSets)-1]))
			if len(colSets) == len(cols) {
				break
			}
		}
	}
	if len(currentWord) > 0 {
		wordGroups = append(wordGroups, currentWord)
	}

	for _, group := range wordGroups {
		p.AddConstraint(constraints.NewValidWordConstraint(p, group, t, alphabet))
	}
	for i := 0; i < numCols; i++ {
		p.AddConstraint(constraints.NewSetCountConstraint(cols[i], colSets[i], isCovering[i]))
		allowedVals := puzzle.ValueSet{}
		for val := range colSets[i] {
			allowedVals[val] = true
		}
		for _, v := range cols[i] {
			newVals := puzzle.ValueSet{}
			for val := range v.Values() {
				if allowedVals[val] {
					newVals[val] = true
				}
			}
			v.SetValueRange(newVals)
		}
	}

	return &Dropquote{
		variables: allVars,
	}
}

func (d *Dropquote) String() string {
	result := ""
	for _, v := range d.variables {
		val := v.Values().Value()
		if val != nil {
			result += val.Str()
		} else {
			result += " "
		}
	}
	return result
}
