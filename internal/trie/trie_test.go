package trie

import (
	"os"
	"testing"
)

var dict = []string{
	"a",
	"at",
	"the",
	"then",
	"boggle",
	"cat",
	"example",
	"exam",
	"exams",
	"examples",
}

func initTrie() *Trie {
	t := New()
	for _, word := range dict {
		t.Add(word)
	}
	return t
}

func initTrieFromFile() *Trie {
	t := New()
	tempFile, err := os.CreateTemp(os.TempDir(), "test")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tempFile.Name())
	for _, word := range dict {
		tempFile.WriteString(word + "\n")
	}
	tempFile.Close()
	t.AddFile(tempFile.Name())
	return t
}

var prefixTests = []struct {
	input string
	want  bool
}{
	{"a", true},
	{"at", true},
	{"att", false},
	{"the", true},
	{"t", true},
	{"", true},
	{"examp", true},
	{"z", false},
}

func TestHasPrefix(t *testing.T) {
	tr := initTrie()
	for _, tt := range prefixTests {
		got := tr.HasPrefix(tt.input)
		if got != tt.want {
			t.Fatalf("input: %s, got: %t, want: %t\n", tt.input, got, tt.want)
		}
	}
}

var wordTests = []struct {
	input string
	want  bool
}{
	{"a", true},
	{"at", true},
	{"att", false},
	{"the", true},
	{"then", true},
	{"t", false},
	{"", false},
	{"examp", false},
	{"z", false},
}

func TestHasWord(t *testing.T) {
	tr := initTrie()
	for _, tt := range wordTests {
		got := tr.HasWord(tt.input)
		if got != tt.want {
			t.Fatalf("input: %s, got: %t, want: %t\n", tt.input, got, tt.want)
		}
	}
}

func TestAddFile(t *testing.T) {
	tr := initTrieFromFile()
	for _, tt := range wordTests {
		got := tr.HasWord(tt.input)
		if got != tt.want {
			t.Fatalf("input: %s, got: %t, want: %t\n", tt.input, got, tt.want)
		}
	}
}
