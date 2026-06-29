# Puzzle API

This directory contains the higher-level API for modeling puzzles and constraint satisfaction problems. It serves as an abstraction layer over the lower-level `csp` solver engine. While the `csp` package relies on integer representations and raw decision indices, the `puzzle` package allows you to work with typed variables and arbitrary value sets (such as strings, booleans, and custom objects) to construct puzzles in a more intuitive and domain-specific way.

## Files

### `puzzle.go`
This is the core file defining the `Puzzle` API, and the only file in this directory. It introduces several types to bridge high-level puzzle definitions with the underlying integer-based solver:

- **`Value` and `ValueSet`**: Provide a wrapper around arbitrary data types (like `int`, `bool`, or `string`) so they can be assigned unique integer IDs. This mapping makes it possible to define a puzzle's domain using real-world values while keeping the solver efficient.
- **`Variable`**: Represents an individual unknown in the puzzle, tracking its currently possible `ValueSet`. Variables can be restricted to specific values or linked to other variables using `MarkEqual`, which merges them via a disjoint-set data structure to reduce the overall search space.
- **`Constraint`**: A high-level interface for rules that apply to a specific group of `Variable`s.
- **`Puzzle`**: The central struct where variables, values, and constraints are registered. When its `Solve()` method is called, it translates the high-level `Variable`s and `Constraint`s into `csp.Decision`s and `csp.ConstraintChecker`s (using internal components like `constraintShim` and `valueSetConstraint`) and invokes the core `csp.Problem` solver.
