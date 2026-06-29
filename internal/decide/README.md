# Decide

This directory contains heuristics and strategy implementations for the constraint satisfaction problem (CSP) solver's decision-making process. When the `csp` solver is searching for a solution, it must recursively pick an undetermined variable (a `csp.Decision`) and try restricting it to one of its possible values. The implementations in this package dictate *which* variable the solver should guess next, which can drastically impact the performance and search space of the solver.

## Files

### `decide.go`
This is the only file in the directory. It provides concrete implementations of the `csp.Decider` interface:

- **`First`**: The simplest, naive strategy. It always selects the very first undecided variable it encounters.
- **`Min`**: Implements the "Minimum Remaining Values" (MRV) heuristic. It iterates through all undecided variables and selects the one with the fewest remaining possibilities. This is a standard and highly effective optimization for CSPs because it minimizes the branching factor of the search tree and quickly surfaces conflicts.
- **`MinMin`**: A more advanced group-based heuristic. It evaluates all constraint groups to find the group with the fewest total remaining options across all its undecided variables. It then selects the variable within that group that has the fewest remaining options. This targets the most highly constrained parts of the puzzle first.
