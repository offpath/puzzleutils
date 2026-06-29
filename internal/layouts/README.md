# Layouts

This directory provides higher-level structural data models (layouts) for organizing `puzzle.Variable`s. Rather than manually managing large, unstructured lists of variables, these layouts allow puzzles to be defined using familiar spatial or relational structures like 2D grids and node-arc graphs. 

## Files

### `graph.go`
Defines the structures needed to model graph-based puzzles (such as Slitherlink or generalized loop puzzles).
- **`Node`, `Arc`, and `Graph`**: Form the foundational building blocks of a graph. Each `Arc` seamlessly embeds a boolean `puzzle.Variable`, which represents whether that specific edge is part of the solution (e.g., whether a line is drawn).
- **`SquareGridGraph`**: A specialized utility layout that automatically constructs an orthogonal 2D grid of nodes and correctly connects adjacent nodes horizontally and vertically with arcs.

### `grid.go`
Defines the `Grid` data structure, which is widely used for cell-based puzzles like Sudoku and Nonograms. 
- It manages a 2D array of `puzzle.Variable`s, abstracting away the math required to translate (row, col) coordinates to a 1D slice.
- It provides convenient helper methods (`GetRow`, `GetCol`, `GetRect`) to easily extract slices of variables when applying constraints (e.g., getting a 3x3 block in Sudoku to apply a uniqueness constraint).
- It includes utilities for initializing a grid from a string pattern and converting the solved grid state back into a human-readable string.
