package main

import (
	"reflect"
	"testing"
)

func TestRuleMatches(t *testing.T) {
	// Rule for letter 'a' in a 5-letter word:
	// 'a' must appear at least once (MinCount: 1), max 2 times.
	// Index 0 must NOT be 'a' (Absent at index 0).
	// Index 2 MUST be 'a' (Present at index 2).
	rule := NewRule('a', 5)
	rule.MinCount = 1
	rule.MaxCount = 2
	rule.Positions[0] = Absent
	rule.Positions[2] = Present

	tests := []struct {
		word string
		want bool
	}{
		{"abaca", false}, // Index 0 is 'a' (forbidden by Positions[0] == Absent)
		{"crane", true},  // c:0, r:1, a:2, n:3, e:4 (Pos[2] == Present satisfied)
		{"plate", true},  // p:0, l:1, a:2, t:3, e:4 (Pos[2] == Present satisfied)
		{"shart", true},  // s:0, h:1, a:2, r:3, t:4
		{"slate", true},  // s:0, l:1, a:2, t:3, e:4
		{"climb", false}, // No 'a' (count=0 < MinCount=1)
		{"aaaaa", false}, // count=5 > MaxCount=2
	}

	for _, tt := range tests {
		got := rule.Matches(tt.word)
		if got != tt.want {
			t.Errorf("rule.Matches(%q) = %v, want %v", tt.word, got, tt.want)
		}
	}
}

func TestExtractRules(t *testing.T) {
	t.Run("Repeated letters in guess and hidden", func(t *testing.T) {
		// hidden = "apple", guess = "puppy"
		rules, err := ExtractRules("apple", "puppy")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		ruleMap := make(map[rune]Rule)
		for _, r := range rules {
			ruleMap[r.Letter] = r
		}

		// 'p' rule
		pRule, ok := ruleMap['p']
		if !ok {
			t.Fatalf("rule for 'p' missing")
		}
		if pRule.MinCount != 2 || pRule.MaxCount != 2 {
			t.Errorf("pRule count = (%d, %d), want (2, 2)", pRule.MinCount, pRule.MaxCount)
		}
		if pRule.Positions[0] != Absent || pRule.Positions[2] != Present || pRule.Positions[3] != Absent {
			t.Errorf("pRule positions = %v, want [0]:Absent, [2]:Present, [3]:Absent", pRule.Positions)
		}

		// 'u' rule
		uRule, ok := ruleMap['u']
		if !ok {
			t.Fatalf("rule for 'u' missing")
		}
		if uRule.MinCount != 0 || uRule.MaxCount != 0 {
			t.Errorf("uRule count = (%d, %d), want (0, 0)", uRule.MinCount, uRule.MaxCount)
		}
	})

	t.Run("Length mismatch error", func(t *testing.T) {
		_, err := ExtractRules("apple", "cat")
		if err == nil {
			t.Errorf("expected error for mismatched word lengths, got nil")
		}
	})
}

func TestMergeRules(t *testing.T) {
	// Guess 1: hidden = "apple", guess = "stare"
	g1Rules, err := ExtractRules("apple", "stare")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Guess 2: hidden = "apple", guess = "puppy"
	g2Rules, err := ExtractRules("apple", "puppy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	merged, err := MergeRules(g1Rules, g2Rules)
	if err != nil {
		t.Fatalf("unexpected error merging rules: %v", err)
	}

	ruleMap := make(map[rune]Rule)
	for _, r := range merged {
		ruleMap[r.Letter] = r
	}

	// 'e' from guess 1 (stare vs apple -> e is Green at index 4)
	eRule, ok := ruleMap['e']
	if !ok || eRule.Positions[4] != Present {
		t.Errorf("expected 'e' to be Present at position 4 in merged rules, got %v", eRule)
	}

	// 'p' from guess 2 (puppy vs apple -> p count is 2, Green at index 2)
	pRule, ok := ruleMap['p']
	if !ok || pRule.MinCount != 2 || pRule.MaxCount != 2 || pRule.Positions[2] != Present {
		t.Errorf("expected 'p' to have count 2 and Present at index 2, got %v", pRule)
	}
}

func TestFilterWords(t *testing.T) {
	candidates := []string{"apple", "apply", "ample", "stare", "paper", "puppy"}
	
	// Extract rules for guess "stare" against hidden "apple"
	rules, err := ExtractRules("apple", "stare")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	filtered := FilterWords(candidates, rules)
	want := []string{"apple", "ample"}

	if !reflect.DeepEqual(filtered, want) {
		t.Errorf("FilterWords() = %v, want %v", filtered, want)
	}
}



