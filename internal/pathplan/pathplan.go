package pathplan
import (
        "container/heap"
        "fmt"
)
type Board interface {
        String() string
        Clone() Board
        IsSolved() bool
        IsFailed() bool
        Heuristic() int
        Controls() []string
        DoControl(c string)
}
type State struct {
        prev *State
        control string
        b Board
        trueDistance int
        heuristicDistance int
	}

type StateHeap []*State
func (s StateHeap) Len() int { return len(s) }
func (s StateHeap) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s StateHeap) Less(i, j int) bool { return s[i].heuristicDistance < s[j].heuristicDistance }
func (s *StateHeap) Push(x interface{}) {
        *s = append(*s, x.(*State))
}
func (s *StateHeap) Pop() interface{} {
        old := *s
        n := len(old)
        x := old[n-1]
        *s = old[0:n-1]
        return x
}

func DijkstraSolve(b Board) []string {
        var s StateHeap
        totalExamined := 0
        maxQueued := 0
        seen := map[string]bool{b.String():true}
        heap.Init(&s)
        heap.Push(&s, &State{nil, "", b, 0, 0})
        for len(s) > 0 {
                totalExamined++
                if len(s) > maxQueued {
                        maxQueued = len(s)
                }
                state := heap.Pop(&s).(*State)
                for _, c := range state.b.Controls() {
                        newB := state.b.Clone()
                        newB.DoControl(c)
                        if str := newB.String(); !seen[str] {
                                seen[str] = true
                                if newB.IsFailed() {

                                       continue
                                } else if newB.IsSolved() {
                                        fmt.Println("Solved!", totalExamined, maxQueued)
                                        result := []string{c}
                                        for cur := state; cur.prev != nil; cur = cur.prev {
                                                result = append(result, cur.control)
                                        }
                                        return result
                                }
                                heap.Push(&s, &State{state, c, newB, state.trueDistance + 1, state.trueDistance + 1})
                        }
                }

        }
        fmt.Println("No solution!")
        return nil
}

func AStarSolve(b Board) []string {
        var s StateHeap
        totalExamined := 0
        maxQueued := 0
        seen := map[string]bool{b.String():true}
        heap.Init(&s)
        heap.Push(&s, &State{nil, "", b, 0, 0})
        for len(s) > 0 {
                totalExamined++
                if len(s) > maxQueued {
                        maxQueued = len(s)
                }
                state := heap.Pop(&s).(*State)
                for _, c := range state.b.Controls() {
                        newB := state.b.Clone()
                        newB.DoControl(c)
                        if str := newB.String(); !seen[str] {
                                seen[str] = true
                                if newB.IsFailed() {
                                        continue
                                } else if newB.IsSolved() {
                                        fmt.Println("Solved!", totalExamined, maxQueued)
                                        result := []string{c}
                                        for cur := state; cur.prev != nil; cur = cur.prev {
                                                result = append(result, cur.control)
                                        }
                                        return result
                                }
                                heap.Push(&s, &State{state, c, newB, state.trueDistance + 1, state.trueDistance + 1 + newB.Heuristic()})
                        }
                }

        }
        fmt.Println("No solution!")
        return nil
}

func IDAStarSolve(b *Board) {

}
