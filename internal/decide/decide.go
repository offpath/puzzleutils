package decide

import (
	"puzzleutils/internal/csp"
)

type First struct {}
func (dec *First) Decide(d []*csp.Decision, g []*csp.Group) *csp.Decision {
	return d[0]
}

type Min struct {}
func (dec *Min) Decide(d []*csp.Decision, g []*csp.Group) *csp.Decision {
	min := d[0]
	for i := 1; i < len(d); i++ {
		if d[i].Count() < min.Count() {
			min = d[i]
		}
	}
	return min
}

type MinMin struct {}
func (dec *MinMin) Decide(decisions []*csp.Decision, groups []*csp.Group) *csp.Decision {
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

