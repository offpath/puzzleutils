package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// extractFiveLetterWords reads the given text file and returns a slice containing all 5-letter words.
func extractFiveLetterWords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if len(word) == 5 {
			words = append(words, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return words, nil
}

func main() {
	filename := "datafiles/ospd2.txt"
	words, err := extractFiveLetterWords(filename)
	if err != nil {
		log.Fatalf("Error loading 5-letter words: %v", err)
	}

	fmt.Printf("Successfully extracted %d 5-letter words from %s.\n", len(words), filename)

	answersFile := "datafiles/wordle_answers.txt"
	answers, err := extractFiveLetterWords(answersFile)
	if err != nil {
		log.Printf("Note: Could not load answer list from %s: %v\n", answersFile, err)
	} else {
		fmt.Printf("Successfully loaded %d original Wordle answer words from %s.\n", len(answers), answersFile)
	}

	// Example: Extract and merge rules across multiple subsequent guesses
	hidden := "apple"
	guesses := []string{"stare", "plans", "puppy"}

	var allGuessRules [][]Rule
	for _, g := range guesses {
		r, err := ExtractRules(hidden, g)
		if err != nil {
			log.Fatalf("Error extracting rules for guess %q: %v", g, err)
		}
		allGuessRules = append(allGuessRules, r)
	}

	mergedRules, err := MergeRules(allGuessRules...)
	if err != nil {
		log.Fatalf("Error merging rules: %v", err)
	}

	fmt.Printf("\nMerged rules after guesses %v against hidden %q:\n", guesses, hidden)
	for _, r := range mergedRules {
		fmt.Println(" ", r)
	}

	// Filter candidate words using merged rules
	remaining := FilterWords(words, mergedRules)
	fmt.Printf("\nWords remaining in dictionary after merged rules (%d candidates): %v\n",
		len(remaining), remaining)

	if len(answers) > 0 {
		remainingAnswers := FilterWords(answers, mergedRules)
		fmt.Printf("Words remaining in answer list after merged rules (%d candidates): %v\n",
			len(remainingAnswers), remainingAnswers)

		fmt.Printf("\nEvaluating %d guess candidates (ospd2.txt) across %d answer words...\n", len(words), len(answers))
		bestStat, topStats := FindOptimalStartingWord(answers, words, words)

		fmt.Printf("\n=========================================\n")
		fmt.Printf(" OPTIMAL STARTING WORD: %s\n", strings.ToUpper(bestStat.Word))
		fmt.Printf(" Total words filtered: %d\n", bestStat.TotalWordsFiltered)
		fmt.Printf(" Average words remaining: %.2f\n", bestStat.AvgWordsRemaining)
		fmt.Printf("=========================================\n")

		fmt.Println("\nTop 10 Starting Words:")
		for i := 0; i < min(10, len(topStats)); i++ {
			st := topStats[i]
			fmt.Printf(" %2d. %-6s | Total Filtered: %d | Avg Remaining: %.2f\n",
				i+1, st.Word, st.TotalWordsFiltered, st.AvgWordsRemaining)
		}

		fmt.Printf("\nEvaluating optimal starting PAIRS of words...\n")
		topK := 200 // Top 200 single starting words paired against all 8,548 candidate guesses
		bestPair, topPairStats := FindOptimalStartingPair(answers, words, words, topK)

		fmt.Printf("\n=========================================\n")
		fmt.Printf(" OPTIMAL STARTING PAIR: (%s, %s)\n", strings.ToUpper(bestPair.Word1), strings.ToUpper(bestPair.Word2))
		fmt.Printf(" Total words filtered: %d\n", bestPair.TotalWordsFiltered)
		fmt.Printf(" Average words remaining: %.2f\n", bestPair.AvgWordsRemaining)
		fmt.Printf("=========================================\n")

		fmt.Println("\nTop 10 Starting Pairs:")
		for i := 0; i < min(10, len(topPairStats)); i++ {
			p := topPairStats[i]
			fmt.Printf(" %2d. (%-6s, %-6s) | Total Filtered: %d | Avg Remaining: %.2f\n",
				i+1, p.Word1, p.Word2, p.TotalWordsFiltered, p.AvgWordsRemaining)
		}

		fmt.Printf("\nEvaluating optimal 3-WORD STARTING OPENINGS...\n")
		topKPairs := 50 // Top 50 starting pairs evaluated with all 8,548 candidate 3rd guesses
		bestTriplet, topTripletStats := FindOptimalStartingTriplet(answers, words, words, topKPairs)

		fmt.Printf("\n=========================================\n")
		fmt.Printf(" OPTIMAL 3-WORD OPENING: (%s, %s, %s)\n",
			strings.ToUpper(bestTriplet.Word1), strings.ToUpper(bestTriplet.Word2), strings.ToUpper(bestTriplet.Word3))
		fmt.Printf(" Total words filtered: %d\n", bestTriplet.TotalWordsFiltered)
		fmt.Printf(" Average words remaining: %.2f\n", bestTriplet.AvgWordsRemaining)
		fmt.Printf("=========================================\n")

		fmt.Println("\nTop 10 3-Word Openings:")
		for i := 0; i < min(10, len(topTripletStats)); i++ {
			tr := topTripletStats[i]
			fmt.Printf(" %2d. (%-6s, %-6s, %-6s) | Total Filtered: %d | Avg Remaining: %.2f\n",
				i+1, tr.Word1, tr.Word2, tr.Word3, tr.TotalWordsFiltered, tr.AvgWordsRemaining)
		}

		fmt.Printf("\nEvaluating all orderings of (%s, %s, %s)...\n",
			bestTriplet.Word1, bestTriplet.Word2, bestTriplet.Word3)
		permStats := EvaluatePermutations(answers, words, []string{bestTriplet.Word1, bestTriplet.Word2, bestTriplet.Word3})
		fmt.Println("\nOrdering Analysis (ranked by cumulative reduction across guesses 1, 2, and 3):")
		for i, ps := range permStats {
			fmt.Printf(" %d. %-24s | After 1: %6.2f | After 2: %5.2f | After 3: %4.2f\n",
				i+1, fmt.Sprintf("(%s, %s, %s)", ps.Order[0], ps.Order[1], ps.Order[2]),
				ps.AvgRemainingAfter1, ps.AvgRemainingAfter2, ps.AvgRemainingAfter3)
		}
	}
}
