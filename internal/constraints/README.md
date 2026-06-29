# Constraints

This directory provides implementations of the `Constraint` interface (defined in the `puzzle` package) to enforce various logical rules over groups of variables. These constraints are used by the higher-level puzzle types (like Sudoku, Nonogram, Dropquote, etc.) to model puzzle-specific mechanics.

## Files

### `buildup.go`
Defines `BuildupSet`, a specialized utility structure used internally by some of the more complex constraints (like `NonogramConstraint` and `ValidWordConstraint`). It tracks valid combinations of values during constraint evaluation and efficiently eliminates impossible values from the variables' domains once all combinations are checked.

### `constraints.go`
Contains generic constraints, most notably the `UniqueConstraint`. This constraint enforces that all variables in a given group must take on distinct, unique values. It is widely used in puzzles like Sudoku.

### `gridgraph.go`
Contains specialized constraints for graph and loop-based logic puzzles, specifically Slitherlink. It defines:
- **`GridGraphPointConstraint`**: Ensures valid line connections at intersections (degree constraints).
- **`GridGraphBoxConstraint`**: Enforces the numbered clues inside grid cells (exact count of surrounding edges).
- **`GridGraphLoopConstraint`**: Prevents the formation of multiple disjoint loops or premature closed loops, ensuring a single continuous path.

### `nonogram.go`
Defines the `NonogramConstraint`, which enforces the run-length encoding clues of Nonogram (Picross) puzzles. It ensures that the sequence of shaded blocks in a row or column exactly matches the provided lengths.

### `setcount.go`
Defines the `SetCountConstraint`. This constraint ensures that specific values appear an exact target number of times across a group of variables. It is used in puzzles where frequency of elements is fixed.

### `validword.go`
Defines the `ValidWordConstraint`, used primarily in word puzzles like Dropquotes and Cryptograms. It takes a sequence of variables and validates that they form a valid word (or prefix) according to a provided `Trie` dictionary, pruning impossible letter assignments.
