# Trie

This directory contains a standard Trie (prefix tree) data structure implementation. The `Trie` is highly optimized for fast prefix and full-word lookups. It is a critical component for any word-based puzzles in this repository (such as Cryptograms and Dropquotes), as it allows the CSP solver to instantly prune paths that result in invalid letter sequences without having to wait until an entire word is fully assigned.

## Files

### `trie.go`
This is the sole non-test file in the directory. It defines the core `Trie` structure and its internal `node` representation. It provides the following key functionalities:

- **`Add(s string)`**: Inserts a new string into the trie. All strings are normalized to uppercase.
- **`HasPrefix(s string) bool`**: Checks if the given string is a valid prefix of *any* word currently stored in the trie. This is the primary method used by the CSP engine (e.g., via `ValidWordConstraint`) to eagerly prune invalid partial assignments.
- **`HasWord(s string) bool`**: Checks if the given string exists exactly as a complete word in the trie.
- **`AddFile(s string)`**: A convenience utility that opens a text file, reads it line by line, and adds every non-empty line into the trie. This is typically used to load large dictionary files (like `ospd2.txt`) into memory at startup.
