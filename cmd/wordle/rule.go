package main

import (
	"fmt"
	"sort"
	"strings"
)

// PositionState defines the status of a letter at a specific index in the hidden word.
type PositionState int

const (
	Unknown PositionState = iota // Letter may or may not be at this position
	Present                      // Letter MUST be at this position (Green)
	Absent                       // Letter CANNOT be at this position (Yellow/Grey)
)

func (ps PositionState) String() string {
	switch ps {
	case Present:
		return "Present"
	case Absent:
		return "Absent"
	default:
		return "Unknown"
	}
}

// Rule represents the constraints learned for a single letter in Wordle.
type Rule struct {
	Letter    rune            // The letter associated with this rule (e.g. 'a')
	MinCount  int             // Minimum number of times this letter appears in the hidden word
	MaxCount  int             // Maximum number of times this letter appears in the hidden word
	Positions []PositionState // Positional constraints for each index (0-based)
}

// NewRule creates a default Rule for a letter given the word length.
func NewRule(letter rune, wordLen int) Rule {
	positions := make([]PositionState, wordLen)
	for i := range positions {
		positions[i] = Unknown
	}
	return Rule{
		Letter:    letter,
		MinCount:  0,
		MaxCount:  wordLen,
		Positions: positions,
	}
}

// Matches checks whether a candidate word satisfies this rule.
func (r Rule) Matches(word string) bool {
	runes := []rune(strings.ToLower(word))
	if len(runes) != len(r.Positions) {
		return false
	}

	count := 0
	for i, ch := range runes {
		if ch == r.Letter {
			count++
			if r.Positions[i] == Absent {
				return false
			}
		} else {
			if r.Positions[i] == Present {
				return false
			}
		}
	}

	if count < r.MinCount || count > r.MaxCount {
		return false
	}

	return true
}

func (r Rule) String() string {
	var posStr []string
	for i, p := range r.Positions {
		if p != Unknown {
			posStr = append(posStr, fmt.Sprintf("[%d]:%s", i, p))
		}
	}
	posInfo := "none"
	if len(posStr) > 0 {
		posInfo = strings.Join(posStr, ", ")
	}
	return fmt.Sprintf("Rule[Letter:'%c', Min:%d, Max:%d, Pos:{%s}]",
		r.Letter, r.MinCount, r.MaxCount, posInfo)
}

type feedbackType int

const (
	feedbackGrey feedbackType = iota
	feedbackYellow
	feedbackGreen
)

// ExtractRules evaluates a guess word against a hidden target word according to Wordle rules
// and returns the learned Rule for each distinct letter present in the guess.
func ExtractRules(hidden, guess string) ([]Rule, error) {
	hiddenRunes := []rune(strings.ToLower(hidden))
	guessRunes := []rune(strings.ToLower(guess))

	if len(hiddenRunes) != len(guessRunes) {
		return nil, fmt.Errorf("hidden and guess words must have the same length (got %d and %d)", len(hiddenRunes), len(guessRunes))
	}

	n := len(guessRunes)
	feedback := make([]feedbackType, n)
	usedHidden := make([]bool, n)

	// First pass: mark Green (exact matches)
	for i := 0; i < n; i++ {
		if guessRunes[i] == hiddenRunes[i] {
			feedback[i] = feedbackGreen
			usedHidden[i] = true
		}
	}

	// Second pass: mark Yellow (present in wrong position) and Grey (absent)
	for i := 0; i < n; i++ {
		if feedback[i] == feedbackGreen {
			continue
		}
		gLetter := guessRunes[i]
		matched := false
		for j := 0; j < n; j++ {
			if !usedHidden[j] && hiddenRunes[j] == gLetter {
				usedHidden[j] = true
				matched = true
				break
			}
		}
		if matched {
			feedback[i] = feedbackYellow
		} else {
			feedback[i] = feedbackGrey
		}
	}

	// Group information by unique letter in guess
	letterIndices := make(map[rune][]int)
	for i, ch := range guessRunes {
		letterIndices[ch] = append(letterIndices[ch], i)
	}

	var rules []Rule
	for ch, indices := range letterIndices {
		rule := NewRule(ch, n)
		matchedCount := 0
		hasGrey := false

		for _, idx := range indices {
			switch feedback[idx] {
			case feedbackGreen:
				rule.Positions[idx] = Present
				matchedCount++
			case feedbackYellow:
				rule.Positions[idx] = Absent
				matchedCount++
			case feedbackGrey:
				rule.Positions[idx] = Absent
				hasGrey = true
			}
		}

		rule.MinCount = matchedCount
		if hasGrey {
			rule.MaxCount = matchedCount
		} else {
			rule.MaxCount = n
		}

		rules = append(rules, rule)
	}

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Letter < rules[j].Letter
	})

	return rules, nil
}

// MergeRule combines two Rules for the same letter into a single merged Rule.
func MergeRule(r1, r2 Rule) (Rule, error) {
	if r1.Letter != r2.Letter {
		return Rule{}, fmt.Errorf("cannot merge rules for different letters: '%c' vs '%c'", r1.Letter, r2.Letter)
	}
	if len(r1.Positions) != len(r2.Positions) {
		return Rule{}, fmt.Errorf("cannot merge rules with different word lengths: %d vs %d", len(r1.Positions), len(r2.Positions))
	}

	merged := NewRule(r1.Letter, len(r1.Positions))
	merged.MinCount = max(r1.MinCount, r2.MinCount)
	merged.MaxCount = min(r1.MaxCount, r2.MaxCount)

	for i := range merged.Positions {
		merged.Positions[i] = mergePosState(r1.Positions[i], r2.Positions[i])
	}

	return merged, nil
}

func mergePosState(s1, s2 PositionState) PositionState {
	if s1 == Present || s2 == Present {
		return Present
	}
	if s1 == Absent || s2 == Absent {
		return Absent
	}
	return Unknown
}

// MergeRules merges multiple slices of Rules (e.g. accumulated from subsequent guesses)
// into a single consolidated slice of Rules.
func MergeRules(ruleSlices ...[]Rule) ([]Rule, error) {
	ruleMap := make(map[rune]Rule)

	for _, slice := range ruleSlices {
		for _, r := range slice {
			existing, found := ruleMap[r.Letter]
			if !found {
				ruleMap[r.Letter] = r
			} else {
				merged, err := MergeRule(existing, r)
				if err != nil {
					return nil, err
				}
				ruleMap[r.Letter] = merged
			}
		}
	}

	mergedRules := make([]Rule, 0, len(ruleMap))
	for _, r := range ruleMap {
		mergedRules = append(mergedRules, r)
	}

	sort.Slice(mergedRules, func(i, j int) bool {
		return mergedRules[i].Letter < mergedRules[j].Letter
	})

	return mergedRules, nil
}

// FilterWords takes a slice of candidate words and a slice of Rules,
// and returns a new slice containing only the words that satisfy all given Rules.
func FilterWords(words []string, rules []Rule) []string {
	var matching []string
	for _, word := range words {
		valid := true
		for _, r := range rules {
			if !r.Matches(word) {
				valid = false
				break
			}
		}
		if valid {
			matching = append(matching, word)
		}
	}
	return matching
}



