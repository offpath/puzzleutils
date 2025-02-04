package trie

import (
	"bufio"
	"os"
	"strings"
)

type Trie struct {
	root *node
}

type node struct {
	word     bool
	children map[rune]*node
}

func New() *Trie {
	return &Trie{&node{}}
}

func (t *Trie) Add(s string) {
	s = strings.ToUpper(s)
	n := t.root
	for _, c := range s {
		if n.children == nil {
			n.children = map[rune]*node{}
		}
		if n.children[c] == nil {
			n.children[c] = &node{}
		}
		n = n.children[c]
	}
	n.word = true
}

func (t *Trie) HasPrefix(s string) bool {
	s = strings.ToUpper(s)
	n := t.root
	for _, c := range s {
		if n.children == nil || n.children[c] == nil {
			return false
		}
		n = n.children[c]
	}
	return n != nil
}

func (t *Trie) HasWord(s string) bool {
	s = strings.ToUpper(s)
	n := t.root
	for _, c := range s {
		if n.children == nil || n.children[c] == nil {
			return false
		}
		n = n.children[c]
	}
	return n != nil && n.word
}

func (t *Trie) AddFile(s string) {
	readFile, err := os.Open(s)
	if err != nil {
		panic(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		line := strings.TrimSpace(fileScanner.Text())
		if len(line) > 0 {
			t.Add(line)
		}
	}
	readFile.Close()
}
