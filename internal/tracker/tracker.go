package tracker

import (
	"fmt"

	"puzzleutils/internal/csp"
)

func PrintEveryN(n int) csp.DecisionTracker {
	return &printEveryN{0, n}
}

func PrintEveryLogN(n int) csp.DecisionTracker {
	return &printEveryLogN{0, n, 1}
}

type printEveryN struct {
	count int
	multiple int
}

func (pr *printEveryN) CaptureDecision(p *csp.Problem) {
	pr.count++
	if pr.count % pr.multiple == 0 {
		fmt.Printf("Count = %d\n", pr.count)
	}
}

type printEveryLogN struct {
	count int
	multiple int
	next int
}

func (pr *printEveryLogN) CaptureDecision(p *csp.Problem) {
	pr.count++
	if pr.count % pr.next == 0 {
		fmt.Printf("Count = %d\n", pr.count)
		if pr.count / pr.next == pr.multiple {
			pr.next = pr.next * pr.multiple
		}
	}
}
