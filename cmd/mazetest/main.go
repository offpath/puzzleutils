package main

import (
	"fmt"

	"github.com/offpath/puzzleutils/internal/maze"
	"github.com/offpath/puzzleutils/internal/pathplan"
)

func main() {
	board := maze.LoadBoard([]string{
		"XXEXX",
		"X..XX",
		"X.X.X",
		"X...X",
		"XXPXX",
		"X...X",
		"X...X",
		"X...X",
		"X...X",
		"X...X",
		"XXXXX",
	})
	result1 := pathplan.DijkstraSolve(board)
	for _, c := range result1 {
		fmt.Println(c)
	}
	result2 := pathplan.AStarSolve(board)
	for _, c := range result2 {
		fmt.Println(c)
	}
}
