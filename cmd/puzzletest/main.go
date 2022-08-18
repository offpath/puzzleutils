package main

import (
	"fmt"

	"puzzleutils/internal/csp"
	"puzzleutils/internal/puzzle"
)

type DecideFirst struct {}
func (dec *DecideFirst) Decide(d []*csp.Decision, g []*csp.Group) *csp.Decision {
	return d[0]
}

type DecideMin struct {}
func (dec *DecideMin) Decide(d []*csp.Decision, g []*csp.Group) *csp.Decision {
	min := d[0]
	for i := 1; i < len(d); i++ {
		if d[i].Count() < min.Count() {
			min = d[i]
		}
	}
	return min
}

type DecideMinMin struct {}
func (dec *DecideMinMin) Decide(decisions []*csp.Decision, groups []*csp.Group) *csp.Decision {
	var result *csp.Decision
	groupMin := -1
	for _, g := range groups {
		count := 0
		var min *csp.Decision
		for _, d := range g.Decisions() {
			if c := d.Count(); c > 1 {
				count += c - 1
				if min == nil || c < min.Count() {
					min = d
				}
			}
		}
		if count > 0 && (groupMin < 0 || count < groupMin) {
			groupMin = count
			result = min
		}
	}
	return result
}

type Printer struct {
	count int
	multiple int
	p *puzzle.GridPuzzle
}

func (s *Printer) MakeDecision(p *csp.Problem) {
	s.count++
	if s.count % s.multiple == 0 {
		fmt.Printf("Count = %d\n", s.count)
	}
}

func (s *Printer) CaptureSolution(p *csp.Problem) {
	s.p.Print()
}

type Settings struct {
	csp.Printer
	csp.Decider
}



func main() {
	fmt.Println("Hello World!!")
	/*
	sudoku.Init(
		"53..7...." +
			"6..195..." +
			".98....6." +
			"8...6...3" +
			"4..8.3..1" +
			"7...2...6" +
			".6....28." +
			"...419..5" +
			"....8..79")
	*/
	
	sudoku0 := puzzle.NewSudokuPuzzle()
	sudoku0.Init(
		"........." +
			".....3.85" +
			"..1.2...." +
			"...5.7..." +
			"..4...1.." +
			".9......." +
			"5......73" +
			"..2.1...." +
			"....4...9")

	p0 := Printer{0, 100000, sudoku0}
	s0 := Settings{&p0, &DecideFirst{}}
	sudoku0.Solve(s0)
	fmt.Printf("Decisions made: %d\n", p0.count)

	sudoku1 := puzzle.NewSudokuPuzzle()
	sudoku1.Init(
		"........." +
			".....3.85" +
			"..1.2...." +
			"...5.7..." +
			"..4...1.." +
			".9......." +
			"5......73" +
			"..2.1...." +
			"....4...9")

	p1 := Printer{0, 100000, sudoku1}
	s1 := Settings{&p1, &DecideMin{}}
	sudoku1.Solve(s1)
	fmt.Printf("Decisions made: %d\n", p1.count)	

	sudoku2 := puzzle.NewSudokuPuzzle()
	sudoku2.Init(
		"........." +
			".....3.85" +
			"..1.2...." +
			"...5.7..." +
			"..4...1.." +
			".9......." +
			"5......73" +
			"..2.1...." +
			"....4...9")

	p2 := Printer{0, 100000, sudoku2}
	s2 := Settings{&p2, &DecideMinMin{}}
	sudoku2.Solve(s2)
	fmt.Printf("Decisions made: %d\n", p2.count)

	rows := [][]int{
		{8,7,5,7},
		{5,4,3,3},
		{3,3,2,3},
		{4,3,2,2},
		{3,3,2,2},
		{3,4,2,2},
		{4,5,2},
		{3,5,1},
		{4,3,2},
		{3,4,2},
		{4,4,2},
		{3,6,2},
		{3,2,3,1},
		{4,3,4,2},
		{3,2,3,2},
		{6,5},
		{4,5},
		{3,3},
		{3,3},
		{1,1},
	}
	cols := [][]int{
		{1},
		{1},
		{2},
		{4},
		{7},
		{9},
		{2,8},
		{1,8},
		{8},
		{1,9},
		{2,7},
		{3,4},
		{6,4},
		{8,5},
		{1,11},
		{1,7},
		{8},
		{1,4,8},
		{6,8},
		{4,7},
		{2,4},
		{1,4},
		{5},
		{1,4},
		{1,5},
		{7},
		{5},
		{3},
		{1},
		{1},
	}
	nonogram := puzzle.NewNonogramPuzzle(rows, cols)
	p3 := Printer{0, 10000, nonogram}
	s3 := Settings{&p3, &DecideFirst{}}
	nonogram.Solve(s3)
	fmt.Printf("Decisions made: %d\n", p3.count)
}
