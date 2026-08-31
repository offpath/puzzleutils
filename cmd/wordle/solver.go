package main

import (
	"runtime"
	"sort"
	"strings"
	"sync"
)

// WordStat stores analysis metrics for a starting guess word.
type WordStat struct {
	Word               string  // The starting guess word
	TotalWordsFiltered int64   // Total candidate words filtered across all answer words
	AvgWordsRemaining  float64 // Average candidate words remaining after this starting guess
}

// CalculatePattern computes the Wordle feedback pattern integer (0..242) for guess vs target.
func CalculatePattern(target, guess string) int {
	t := []rune(strings.ToLower(target))
	g := []rune(strings.ToLower(guess))
	n := len(g)

	color := make([]int, n)
	usedT := make([]bool, n)

	// Pass 1: Green (exact match)
	for i := 0; i < n; i++ {
		if g[i] == t[i] {
			color[i] = 2 // Green
			usedT[i] = true
		}
	}

	// Pass 2: Yellow / Grey
	for i := 0; i < n; i++ {
		if color[i] == 2 {
			continue
		}
		for j := 0; j < n; j++ {
			if !usedT[j] && t[j] == g[i] {
				color[i] = 1 // Yellow
				usedT[j] = true
				break
			}
		}
	}

	// Base-3 encoding into a single integer pattern ID
	pat := 0
	for i := 0; i < n; i++ {
		pat = pat*3 + color[i]
	}
	return pat
}

// FindOptimalStartingWord iterates over all candidate answers and first guess options,
// evaluating the rules/feedback produced and filtering candidate words using FilterWords.
// Returns all ranked WordStats (best starting word first).
func FindOptimalStartingWord(answers, candidateGuesses, candidatePool []string) (WordStat, []WordStat) {
	numGuesses := len(candidateGuesses)
	stats := make([]WordStat, numGuesses)

	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup
	guessChan := make(chan int, numGuesses)

	for i := 0; i < numGuesses; i++ {
		guessChan <- i
	}
	close(guessChan)

	totalCandidates := int64(len(candidatePool))
	totalAnswers := int64(len(answers))
	maxPossibleFiltered := totalAnswers * totalCandidates

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for gIdx := range guessChan {
				g := candidateGuesses[gIdx]

				// Count feedback pattern distribution for guess g across all candidate words in pool
				var bucketCounts [243]int
				for _, cand := range candidatePool {
					pat := CalculatePattern(cand, g)
					bucketCounts[pat]++
				}

				// Sum remaining matching candidate words across all target answer words
				var totalRemaining int64
				for _, ans := range answers {
					pat := CalculatePattern(ans, g)
					totalRemaining += int64(bucketCounts[pat])
				}

				totalFiltered := maxPossibleFiltered - totalRemaining
				avgRemaining := float64(totalRemaining) / float64(totalAnswers)

				stats[gIdx] = WordStat{
					Word:               g,
					TotalWordsFiltered: totalFiltered,
					AvgWordsRemaining:  avgRemaining,
				}
			}
		}()
	}

	wg.Wait()

	// Sort stats by TotalWordsFiltered descending
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].TotalWordsFiltered == stats[j].TotalWordsFiltered {
			return stats[i].Word < stats[j].Word
		}
		return stats[i].TotalWordsFiltered > stats[j].TotalWordsFiltered
	})

	return stats[0], stats
}

// PairStat stores analysis metrics for a pair of starting guess words.
type PairStat struct {
	Word1              string  // The first starting guess word
	Word2              string  // The second starting guess word
	TotalWordsFiltered int64   // Total candidate words filtered across all answer words after both guesses
	AvgWordsRemaining  float64 // Average candidate words remaining after both guesses
}

// FindOptimalStartingPair searches for the optimal pair of starting words.
// It precomputes pattern matrices and evaluates pairs in parallel.
func FindOptimalStartingPair(answers, candidateGuesses, candidatePool []string, topKFirst int) (PairStat, []PairStat) {
	numGuesses := len(candidateGuesses)
	numCandidates := len(candidatePool)
	numAnswers := len(answers)

	// 1. Precompute Pattern Matrices for zero-allocation matrix lookup
	patMatrix := make([][]uint8, numGuesses)
	for gIdx, g := range candidateGuesses {
		patMatrix[gIdx] = make([]uint8, numCandidates)
		for cIdx, c := range candidatePool {
			patMatrix[gIdx][cIdx] = uint8(CalculatePattern(c, g))
		}
	}

	ansPatMatrix := make([][]uint8, numGuesses)
	for gIdx, g := range candidateGuesses {
		ansPatMatrix[gIdx] = make([]uint8, numAnswers)
		for aIdx, a := range answers {
			ansPatMatrix[gIdx][aIdx] = uint8(CalculatePattern(a, g))
		}
	}

	// 2. Rank single starting words to select top first guesses
	_, singleStats := FindOptimalStartingWord(answers, candidateGuesses, candidatePool)
	if topKFirst <= 0 || topKFirst > len(singleStats) {
		topKFirst = len(singleStats)
	}

	guessIdxMap := make(map[string]int, numGuesses)
	for i, g := range candidateGuesses {
		guessIdxMap[g] = i
	}

	firstGuessIndices := make([]int, topKFirst)
	for i := 0; i < topKFirst; i++ {
		firstGuessIndices[i] = guessIdxMap[singleStats[i].Word]
	}

	// 3. Worker queue for parallel pair evaluation
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	workChan := make(chan int, topKFirst)
	for _, g1Idx := range firstGuessIndices {
		workChan <- g1Idx
	}
	close(workChan)

	var mu sync.Mutex
	var allPairStats []PairStat

	totalCandidates := int64(numCandidates)
	totalAnswers := int64(numAnswers)
	maxPossibleFiltered := totalAnswers * totalCandidates

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var jointBucketCounts [59049]int32
			localStats := make([]PairStat, 0, numGuesses)

			for g1Idx := range workChan {
				g1Name := candidateGuesses[g1Idx]
				g1Pats := patMatrix[g1Idx]
				g1AnsPats := ansPatMatrix[g1Idx]

				for g2Idx := 0; g2Idx < numGuesses; g2Idx++ {
					if g1Idx == g2Idx {
						continue
					}
					g2Pats := patMatrix[g2Idx]
					g2AnsPats := ansPatMatrix[g2Idx]

					for i := range jointBucketCounts {
						jointBucketCounts[i] = 0
					}

					for cIdx := 0; cIdx < numCandidates; cIdx++ {
						patIdx := int(g1Pats[cIdx])*243 + int(g2Pats[cIdx])
						jointBucketCounts[patIdx]++
					}

					var totalRemaining int64
					for aIdx := 0; aIdx < numAnswers; aIdx++ {
						patIdx := int(g1AnsPats[aIdx])*243 + int(g2AnsPats[aIdx])
						totalRemaining += int64(jointBucketCounts[patIdx])
					}

					totalFiltered := maxPossibleFiltered - totalRemaining
					avgRemaining := float64(totalRemaining) / float64(totalAnswers)

					localStats = append(localStats, PairStat{
						Word1:              g1Name,
						Word2:              candidateGuesses[g2Idx],
						TotalWordsFiltered: totalFiltered,
						AvgWordsRemaining:  avgRemaining,
					})
				}
			}

			mu.Lock()
			allPairStats = append(allPairStats, localStats...)
			mu.Unlock()
		}()
	}

	wg.Wait()

	sort.Slice(allPairStats, func(i, j int) bool {
		if allPairStats[i].TotalWordsFiltered == allPairStats[j].TotalWordsFiltered {
			if allPairStats[i].Word1 == allPairStats[j].Word1 {
				return allPairStats[i].Word2 < allPairStats[j].Word2
			}
			return allPairStats[i].Word1 < allPairStats[j].Word1
		}
		return allPairStats[i].TotalWordsFiltered > allPairStats[j].TotalWordsFiltered
	})

	return allPairStats[0], allPairStats
}

// TripletStat stores analysis metrics for a 3-word starting opening.
type TripletStat struct {
	Word1              string  // The first starting guess word
	Word2              string  // The second starting guess word
	Word3              string  // The third starting guess word
	TotalWordsFiltered int64   // Total candidate words filtered across all answer words after 3 guesses
	AvgWordsRemaining  float64 // Average candidate words remaining after 3 guesses
}

// FindOptimalStartingTriplet evaluates 3-word starting openings by taking top pairs
// and finding the best complementary 3rd guess word.
func FindOptimalStartingTriplet(answers, candidateGuesses, candidatePool []string, topKPairs int) (TripletStat, []TripletStat) {
	numGuesses := len(candidateGuesses)
	numCandidates := len(candidatePool)
	numAnswers := len(answers)

	// Precompute Pattern Matrices
	patMatrix := make([][]uint8, numGuesses)
	for gIdx, g := range candidateGuesses {
		patMatrix[gIdx] = make([]uint8, numCandidates)
		for cIdx, c := range candidatePool {
			patMatrix[gIdx][cIdx] = uint8(CalculatePattern(c, g))
		}
	}

	ansPatMatrix := make([][]uint8, numGuesses)
	for gIdx, g := range candidateGuesses {
		ansPatMatrix[gIdx] = make([]uint8, numAnswers)
		for aIdx, a := range answers {
			ansPatMatrix[gIdx][aIdx] = uint8(CalculatePattern(a, g))
		}
	}

	guessIdxMap := make(map[string]int, numGuesses)
	for i, g := range candidateGuesses {
		guessIdxMap[g] = i
	}

	// Rank pairs first to select top starting pairs
	_, topPairStats := FindOptimalStartingPair(answers, candidateGuesses, candidatePool, 100)
	if topKPairs <= 0 || topKPairs > len(topPairStats) {
		topKPairs = len(topPairStats)
	}

	type pairIdx struct {
		g1Idx int
		g2Idx int
		word1 string
		word2 string
	}
	topPairs := make([]pairIdx, topKPairs)
	for i := 0; i < topKPairs; i++ {
		topPairs[i] = pairIdx{
			g1Idx: guessIdxMap[topPairStats[i].Word1],
			g2Idx: guessIdxMap[topPairStats[i].Word2],
			word1: topPairStats[i].Word1,
			word2: topPairStats[i].Word2,
		}
	}

	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	workChan := make(chan pairIdx, topKPairs)
	for _, p := range topPairs {
		workChan <- p
	}
	close(workChan)

	var mu sync.Mutex
	var allTripletStats []TripletStat

	totalCandidates := int64(numCandidates)
	totalAnswers := int64(numAnswers)
	maxPossibleFiltered := totalAnswers * totalCandidates

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localStats := make([]TripletStat, 0, numGuesses)

			for pair := range workChan {
				g1Idx := pair.g1Idx
				g2Idx := pair.g2Idx

				g1Pats := patMatrix[g1Idx]
				g2Pats := patMatrix[g2Idx]
				g1AnsPats := ansPatMatrix[g1Idx]
				g2AnsPats := ansPatMatrix[g2Idx]

				// Pre-calculate remaining candidates list for each answer
				ansRemList := make([][]int, numAnswers)
				for aIdx := 0; aIdx < numAnswers; aIdx++ {
					targetP1 := g1AnsPats[aIdx]
					targetP2 := g2AnsPats[aIdx]
					var rem []int
					for cIdx := 0; cIdx < numCandidates; cIdx++ {
						if g1Pats[cIdx] == targetP1 && g2Pats[cIdx] == targetP2 {
							rem = append(rem, cIdx)
						}
					}
					ansRemList[aIdx] = rem
				}

				// Evaluate each candidate 3rd guess g3
				for g3Idx := 0; g3Idx < numGuesses; g3Idx++ {
					if g3Idx == g1Idx || g3Idx == g2Idx {
						continue
					}

					g3Pats := patMatrix[g3Idx]
					g3AnsPats := ansPatMatrix[g3Idx]

					var totalRemaining int64
					for aIdx := 0; aIdx < numAnswers; aIdx++ {
						targetP3 := g3AnsPats[aIdx]
						rem := ansRemList[aIdx]
						count := 0
						for _, cIdx := range rem {
							if g3Pats[cIdx] == targetP3 {
								count++
							}
						}
						totalRemaining += int64(count)
					}

					totalFiltered := maxPossibleFiltered - totalRemaining
					avgRemaining := float64(totalRemaining) / float64(totalAnswers)

					localStats = append(localStats, TripletStat{
						Word1:              pair.word1,
						Word2:              pair.word2,
						Word3:              candidateGuesses[g3Idx],
						TotalWordsFiltered: totalFiltered,
						AvgWordsRemaining:  avgRemaining,
					})
				}
			}

			mu.Lock()
			allTripletStats = append(allTripletStats, localStats...)
			mu.Unlock()
		}()
	}

	wg.Wait()

	sort.Slice(allTripletStats, func(i, j int) bool {
		if allTripletStats[i].TotalWordsFiltered == allTripletStats[j].TotalWordsFiltered {
			if allTripletStats[i].Word1 == allTripletStats[j].Word1 {
				if allTripletStats[i].Word2 == allTripletStats[j].Word2 {
					return allTripletStats[i].Word3 < allTripletStats[j].Word3
				}
				return allTripletStats[i].Word2 < allTripletStats[j].Word2
			}
			return allTripletStats[i].Word1 < allTripletStats[j].Word1
		}
		return allTripletStats[i].TotalWordsFiltered > allTripletStats[j].TotalWordsFiltered
	})

	return allTripletStats[0], allTripletStats
}

// PermutationStat tracks step-by-step candidate reduction for a 3-word sequence.
type PermutationStat struct {
	Order              []string
	AvgRemainingAfter1 float64
	AvgRemainingAfter2 float64
	AvgRemainingAfter3 float64
	CumulativeAvgSum   float64
}

// EvaluatePermutations evaluates all 6 orderings of a 3-word opening triad.
func EvaluatePermutations(answers, candidatePool []string, triad []string) []PermutationStat {
	perms := [][]string{
		{triad[0], triad[1], triad[2]},
		{triad[0], triad[2], triad[1]},
		{triad[1], triad[0], triad[2]},
		{triad[1], triad[2], triad[0]},
		{triad[2], triad[0], triad[1]},
		{triad[2], triad[1], triad[0]},
	}

	numCandidates := len(candidatePool)
	numAnswers := len(answers)
	totalAnswers := float64(numAnswers)

	var stats []PermutationStat

	for _, p := range perms {
		w1, w2, w3 := p[0], p[1], p[2]

		p1Cand := make([]uint8, numCandidates)
		p1Ans := make([]uint8, numAnswers)
		p2Cand := make([]uint8, numCandidates)
		p2Ans := make([]uint8, numAnswers)
		p3Cand := make([]uint8, numCandidates)
		p3Ans := make([]uint8, numAnswers)

		for cIdx, c := range candidatePool {
			p1Cand[cIdx] = uint8(CalculatePattern(c, w1))
			p2Cand[cIdx] = uint8(CalculatePattern(c, w2))
			p3Cand[cIdx] = uint8(CalculatePattern(c, w3))
		}
		for aIdx, a := range answers {
			p1Ans[aIdx] = uint8(CalculatePattern(a, w1))
			p2Ans[aIdx] = uint8(CalculatePattern(a, w2))
			p3Ans[aIdx] = uint8(CalculatePattern(a, w3))
		}

		// Step 1
		var buckets1 [243]int
		for _, pat := range p1Cand {
			buckets1[pat]++
		}
		var rem1 int64
		for _, pat := range p1Ans {
			rem1 += int64(buckets1[pat])
		}
		avg1 := float64(rem1) / totalAnswers

		// Step 2
		var buckets12 [59049]int
		for cIdx := 0; cIdx < numCandidates; cIdx++ {
			idx := int(p1Cand[cIdx])*243 + int(p2Cand[cIdx])
			buckets12[idx]++
		}
		var rem2 int64
		for aIdx := 0; aIdx < numAnswers; aIdx++ {
			idx := int(p1Ans[aIdx])*243 + int(p2Ans[aIdx])
			rem2 += int64(buckets12[idx])
		}
		avg2 := float64(rem2) / totalAnswers

		// Step 3
		var rem3 int64
		for aIdx := 0; aIdx < numAnswers; aIdx++ {
			t1 := p1Ans[aIdx]
			t2 := p2Ans[aIdx]
			t3 := p3Ans[aIdx]
			count := 0
			for cIdx := 0; cIdx < numCandidates; cIdx++ {
				if p1Cand[cIdx] == t1 && p2Cand[cIdx] == t2 && p3Cand[cIdx] == t3 {
					count++
				}
			}
			rem3 += int64(count)
		}
		avg3 := float64(rem3) / totalAnswers

		stats = append(stats, PermutationStat{
			Order:              []string{w1, w2, w3},
			AvgRemainingAfter1: avg1,
			AvgRemainingAfter2: avg2,
			AvgRemainingAfter3: avg3,
			CumulativeAvgSum:   avg1 + avg2 + avg3,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].CumulativeAvgSum < stats[j].CumulativeAvgSum
	})

	return stats
}



