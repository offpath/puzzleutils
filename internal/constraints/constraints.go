package constraints

import (
	"github.com/offpath/puzzleutils/internal/csp"
	"github.com/offpath/puzzleutils/internal/trie"
)

func Unique(isCovering bool) csp.ConstraintChecker {
	return unique{isCovering}
}

func UniqueCovering() csp.ConstraintChecker {
	return Unique(true)
}

func Equal() csp.ConstraintChecker {
	return equal{}
}

func Set(s map[int]bool) csp.ConstraintChecker {
	return set{s}
}

func SetCount(s map[int]int, isCovering bool) csp.ConstraintChecker {
	return setCount{s, isCovering}
}

func ValidWord(t *trie.Trie, valueSet []string) csp.ConstraintChecker {
	return validWord{t, valueSet}
}

type unique struct {
	isCovering bool
}

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
	// Break early if we're not checking for a covering set.
	if !c.isCovering {
		return true
	}
	// If this is a covering set, there are exactly as many values as decisions.
	for i := 0; i < len(all); i++ {
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

type setCount struct {
	s          map[int]int
	isCovering bool
}

func (c setCount) Init(all []*csp.Decision, size int) {
	s := map[int]bool{}
	for item := range c.s {
		s[item] = true
	}
	for _, d := range all {
		d.RestrictToSet(s)
	}
}

func (c setCount) Apply(all, dirty []*csp.Decision) bool {
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
		}
		if count > target || (c.isCovering && possibleCount < target) {
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
	return true
}

type validWord struct {
	t        *trie.Trie
	valueSet []string
}

func (c validWord) Init(all []*csp.Decision, size int) {}

func (c validWord) Apply(all, dirty []*csp.Decision) bool {
	b := NewBuildupSet(len(all))
	var f func(ds []*csp.Decision, prefix string)
	f = func(ds []*csp.Decision, prefix string) {
		if len(ds) == 0 {
			return
		}
		for i, val := range c.valueSet {
			if !ds[0].Possible(i) {
				continue
			}
			if (len(ds) == 1 && c.t.HasWord(prefix+val)) || (len(ds) > 1 && c.t.HasPrefix(prefix+val)) {
				b.Push(i)
				f(ds[1:], prefix+val)
				b.Pop()
			}
		}
	}
	f(all, "")
	b.Export(all)
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
