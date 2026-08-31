package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExtractFiveLetterWords(t *testing.T) {
	content := "cat\naahed\naalii\nsomewordtoolong\naargh\ndog\n"
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_words.txt")

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp test file: %v", err)
	}

	words, err := extractFiveLetterWords(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"aahed", "aalii", "aargh"}
	if !reflect.DeepEqual(words, want) {
		t.Errorf("extractFiveLetterWords() = %v, want %v", words, want)
	}
}
