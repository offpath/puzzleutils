package main

import (
	"testing"
)

func TestCalculatePatternEquivalence(t *testing.T) {
	hidden := "apple"
	guess := "puppy"
	candidates := []string{"apple", "appal", "appel", "apply", "ample", "stare", "paper", "puppy"}

	// 1. Calculate using ExtractRules & FilterWords
	rules, err := ExtractRules(hidden, guess)
	if err != nil {
		t.Fatalf("ExtractRules error: %v", err)
	}
	filtered := FilterWords(candidates, rules)

	// 2. Calculate using CalculatePattern matching
	targetPat := CalculatePattern(hidden, guess)
	var patternMatched []string
	for _, c := range candidates {
		if CalculatePattern(c, guess) == targetPat {
			patternMatched = append(patternMatched, c)
		}
	}

	if len(filtered) != len(patternMatched) {
		t.Fatalf("Mismatch in count: FilterWords = %d, CalculatePattern = %d", len(filtered), len(patternMatched))
	}
	for i := range filtered {
		if filtered[i] != patternMatched[i] {
			t.Errorf("Mismatch at %d: FilterWords = %s, CalculatePattern = %s", i, filtered[i], patternMatched[i])
		}
	}
}

func TestFindOptimalStartingWordSmall(t *testing.T) {
	answers := []string{"apple", "stare", "crane"}
	guesses := []string{"stare", "crane", "apple", "puppy"}

	best, allStats := FindOptimalStartingWord(answers, guesses, guesses)
	if len(allStats) != len(guesses) {
		t.Errorf("expected %d stats, got %d", len(guesses), len(allStats))
	}
	if best.Word == "" {
		t.Errorf("best word should not be empty")
	}
}

func TestFindOptimalStartingPairSmall(t *testing.T) {
	answers := []string{"apple", "stare", "crane"}
	guesses := []string{"stare", "crane", "apple", "puppy"}

	bestPair, allPairStats := FindOptimalStartingPair(answers, guesses, guesses, 2)
	if len(allPairStats) == 0 {
		t.Fatalf("expected pair stats, got empty slice")
	}
	if bestPair.Word1 == "" || bestPair.Word2 == "" {
		t.Errorf("best pair words should not be empty: %v", bestPair)
	}
}

func TestFindOptimalStartingTripletSmall(t *testing.T) {
	answers := []string{"apple", "stare", "crane"}
	guesses := []string{"stare", "crane", "apple", "puppy"}

	bestTriplet, allStats := FindOptimalStartingTriplet(answers, guesses, guesses, 2)
	if len(allStats) == 0 {
		t.Fatalf("expected triplet stats, got empty slice")
	}
	if bestTriplet.Word1 == "" || bestTriplet.Word2 == "" || bestTriplet.Word3 == "" {
		t.Errorf("best triplet words should not be empty: %v", bestTriplet)
	}
}

func TestEvaluatePermutations(t *testing.T) {
	answers := []string{"apple", "stare", "crane"}
	guesses := []string{"stare", "crane", "apple", "puppy"}

	stats := EvaluatePermutations(answers, guesses, []string{"rails", "conte", "dumpy"})
	if len(stats) != 6 {
		t.Fatalf("expected 6 permutation stats, got %d", len(stats))
	}
}



