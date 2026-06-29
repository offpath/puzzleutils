# CSP (Constraint Satisfaction Problem) Solver

This directory contains the core engine for solving constraint satisfaction problems (CSPs). A CSP consists of a set of decisions (variables) that must be assigned values such that a given set of constraints is satisfied. This package provides a generic, recursive backtracking solver that can be used to model and solve various logic puzzles (like Sudoku, Nonograms, Cryptograms, etc.) or other combinatorial problems.

## Files

### `csp.go`
This is the main engine of the solver and the only non-test file in this directory. It defines the core data structures and logic required to build and solve a CSP:

- **`Problem`**: The primary struct that captures all decisions, groups (constraints), and ephemeral state during the solving process. It provides the `Solve` method which recursively backtracks to find a solution.
- **`Decision`**: Represents a single variable to be determined (e.g., a single cell in a grid). It tracks the remaining possibilities and manages state updates when possibilities are restricted.
- **`Group`**: Represents a collection of `Decision`s bound by a specific constraint.
- **`ConstraintChecker`**: An interface for defining rules (constraints) over a group of decisions. The `Apply` method is used to prune possibilities that violate the constraint.
- **`Decider`**: An interface used to determine the next best `Decision` to attempt during the recursive solve process (e.g., picking the decision with the fewest remaining options).
- **`DecisionTracker` & `SolutionTracker`**: Interfaces used to monitor the progress of the solver and capture valid solutions as they are found.

The algorithm relies on maintaining a dirty set of decisions, propagating constraints via the `ConstraintChecker`, and using a backtracking stack to undo restrictions when a conflict is reached.
