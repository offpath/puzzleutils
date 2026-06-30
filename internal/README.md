# Internal Packages

This directory contains the core logic and engine for modeling and solving various types of logic and combinatorial puzzles. The architecture is modular, separating the low-level solving algorithms from the high-level puzzle definitions, structural layouts, and specific logical constraints.

## Subdirectory Summary

### `csp`
Contains the core recursive backtracking constraint satisfaction problem (CSP) solver engine. It provides the low-level integer-based variables (Decisions) and mechanisms for restricting domains, applying constraints, and backtracking upon conflicts.

### `puzzle`
Provides a higher-level API abstraction over the `csp` solver. It allows puzzles to be defined using domain-specific, typed variables (booleans, strings, custom objects) rather than raw integer indices. It handles mapping these arbitrary values to the underlying solver engine seamlessly.

### `constraints`
Contains concrete implementations of the `puzzle.Constraint` interface. These are reusable logical rules—such as uniqueness, exact occurrences (`SetCount`), valid dictionary words, or graph looping rules—that can be applied to groups of variables to model specific puzzle mechanics.

### `layouts`
Provides higher-level structural data models for organizing `puzzle.Variable`s. Instead of managing flat arrays of variables, layouts allow puzzles to be modeled using familiar relational structures like 2D grids (`Grid`) and node-arc graphs (`SquareGridGraph`), exposing convenient methods to extract rows, columns, or surrounding edges.

### `puzzletypes`
Contains the concrete implementations for specific, playable logic puzzles (like Sudoku, Nonogram, Slitherlink, Cryptogram, Dropquote, and Logic Grid puzzles). These packages act as the "glue", utilizing `layouts` for structure, `constraints` for rules, and the `puzzle` API for variable management to build fully solvable models.

### `decide`
Contains heuristics and strategies for the CSP solver's variable selection process. Implementations like `Min` (Minimum Remaining Values) dictate which undecided variable the solver should guess next, which is crucial for optimizing the solver's search tree and performance.

### `tracker`
Provides utilities for monitoring the real-time progress of the `csp` solver. Because complex puzzles can take varying amounts of time to solve, trackers (like `PrintEveryN` or `PrintEveryLogN`) allow developers to log the number of decisions being made without halting the process.

### `trie`
A highly optimized prefix tree data structure. It is used heavily by the `ValidWordConstraint` in word-based puzzles (like Cryptograms) to quickly verify valid prefixes and prune impossible letter combinations early in the solving process.
