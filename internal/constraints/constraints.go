package constraints

import (
	_ "fmt"

	"github.com/offpath/puzzleutils/internal/csp"
)

func Unique() csp.ConstraintChecker {
	return unique{}
}

func UniqueCovering() csp.ConstraintChecker {
	return &uniqueCovering{}
}

func Equal() csp.ConstraintChecker {
	return equal{}
}

func Set(s map[int]bool) csp.ConstraintChecker {
	return set{s}
}

func SetCountCovering(s map[int]int) csp.ConstraintChecker {
	return setCountCovering{s}
}

type unique struct{}

func (c unique) Init(all []*csp.Decision, size int) {}
func (c unique) Apply(all, dirty []*csp.Decision) bool {
	for _, d := range dirty {
		if v := d.Value(); v >= 0 {
			for _, d2 := range all {
				if d2 != d {
					d2.Restrict(v)
				}
			}
		}
	}
	return true
}

type uniqueCovering struct{ size int }

func (c *uniqueCovering) Init(all []*csp.Decision, size int) {
	c.size = size
}
func (c *uniqueCovering) Apply(all, dirty []*csp.Decision) bool {
	for _, d := range dirty {
		if v := d.Value(); v >= 0 {
			for _, d2 := range all {
				if d2 != d {
					d2.Restrict(v)
				}
			}
		}
	}
	for i := 0; i < c.size; i++ {
		found := false
		for _, d := range all {
			if d.Possible(i) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

type equal struct{}

func (c equal) Init(all []*csp.Decision, size int) {}
func (c equal) Apply(all, dirty []*csp.Decision) bool {
	for _, d := range dirty {
		for _, d2 := range all {
			if d2 != d {
				d2.RestrictToEqual(d)
			}
		}
	}
	return true
}

type set struct {
	s map[int]bool
}

func (c set) Init(all []*csp.Decision, size int) {
	for _, d := range all {
		d.RestrictToSet(c.s)
	}
}
func (c set) Apply(all, dirty []*csp.Decision) bool {
	return true
}

// TODO(dneal): Reimplement uniqueCovering using countCovering to test.
type setCountCovering struct {
	s map[int]int
}

func (c setCountCovering) Init(all []*csp.Decision, size int) {
	var s map[int]bool
	for _, item := range c.s {
		s[item] = true
	}
	for _, d := range all {
		d.RestrictToSet(s)
	}
}

func (c setCountCovering) Apply(all, dirty []*csp.Decision) bool {
	for item, target := range c.s {
		count := 0
		possibleCount := 0
		for _, d := range all {
			if d.Possible(item) {
				possibleCount++
			}
			if d.Value() == item {
				count++
			}
			if possibleCount < target || count > target {
				return false
			}
			if count == target {
				for _, d2 := range all {
					if d2.Value() != item {
						d2.Restrict(item)
					}
				}
			}

		}
	}
	return true
}

type BuildupSet struct {
	size, cursor int
	values       []int
	possibleSets []map[int]bool
}

func NewBuildupSet(size int) *BuildupSet {
	result := &BuildupSet{
		size,
		0,
		make([]int, size),
		make([]map[int]bool, size),
	}
	for i := 0; i < size; i++ {
		result.possibleSets[i] = map[int]bool{}
	}
	return result
}

func (b *BuildupSet) Push(val int) {
	b.values[b.cursor] = val
	b.cursor++
	if b.cursor == b.size {
		for i := 0; i < b.size; i++ {
			b.possibleSets[i][b.values[i]] = true
		}
	}
}

func (b *BuildupSet) Pop() {
	b.cursor--
}

func (b *BuildupSet) Export(decisions []*csp.Decision) {
	for i, s := range b.possibleSets {
		decisions[i].RestrictToSet(s)
	}
}
