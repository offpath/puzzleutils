# Puzzle Types

This directory contains the concrete implementations for various types of logic puzzles. These implementations act as the "glue" that binds together the core CSP engine, the domain-specific `puzzle.Variable` definitions, the structural `layouts`, and the specialized `constraints` to model full, solvable puzzles.

## Files

### `cryptogram.go`
Models Cryptogram puzzles. It assigns a `puzzle.Variable` to each unknown cipher character, maps them to an alphabet domain, and enforces two primary constraints: that the sequence of decoded characters forms valid words (using a dictionary Trie), and that the cipher-to-plaintext character mapping is strictly one-to-one (using a `UniqueConstraint`).

### `dropquote.go`
Models Dropquote puzzles. It processes a grid with empty blanks and columns of "dropped" letters. It ensures that the letters dropping into each column match the exact counts of available letters in the bank (using `SetCountConstraint`), and that the resulting horizontal sequences form valid words (using `ValidWordConstraint`).

### `logic.go`
Models classic grid-based Logic Puzzles (e.g., "Zebra puzzles"). It provides a custom parser and Abstract Syntax Tree (AST) evaluator for puzzle clues. Clues like "X is greater than Y" or "A is not B" are parsed and translated into constraints across different categories of variables. 

### `nonogram.go`
Models Nonogram (Picross) puzzles. It utilizes the 2D `Grid` layout and applies the `NonogramConstraint` to each row and column, enforcing that the contiguous blocks of shaded squares match the provided numerical clues exactly.

### `slitherlink.go`
Models Slitherlink puzzles. It utilizes a `SquareGridGraph` layout where arcs represent possible lines. It enforces point constraints (lines must form continuous paths through intersections), box constraints (numbered cells must have exactly that many lines around them), and loop constraints (lines must form a single, continuous, non-intersecting closed loop).

### `sudoku.go`
Models standard 9x9 Sudoku puzzles. It builds a `Grid` layout and registers a `UniqueConstraint` on every row, column, and 3x3 sub-grid, ensuring that the numbers 1 through 9 appear exactly once in each grouping.
