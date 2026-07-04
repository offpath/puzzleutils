# Multiple Solutions Design

## The Problem

The current puzzle solving engine is designed to halt and return as soon as it finds the first valid solution to a given puzzle. For strictly constrained logic puzzles, such as a well-formed Sudoku or Nonogram, this behavior is ideal because there is only one unique solution to find.

However, for more open-ended puzzles—particularly word puzzles like Cryptograms or Dropquotes—a puzzle might be slightly underconstrained and possess multiple technically valid solutions. When this happens, the solver might immediately return the first mathematical configuration it finds. Unfortunately, this first solution might rely on highly obscure, strange, or archaic words from the dictionary, ignoring a much more natural and common human-readable solution that exists elsewhere in the search space. 

Currently, the solver provides no mechanism to discover these alternative solutions. Developers need a way to instruct the engine to not stop at the first success, but to continue exploring the remaining possibilities so they can gather, evaluate, and choose the most desirable solution from all valid options.

## The Technical Plan

To support finding multiple solutions, we will introduce a straightforward callback mechanism that lets developers hook into the solver's lifecycle:

1. **Solution Callback Interface**: We will update the underlying puzzle engine to accept a callback function (a `SolutionTracker`) that fires every time a valid solution is found. 
2. **Continue or Stop**: This callback will return a simple "yes/no" signal indicating whether the solver should keep searching. If the callback says "stop," the solver finishes immediately. If the callback says "continue," the solver will cleverly treat the found solution as if it were a dead end, backtrack, and continue exploring the remaining possibilities.
3. **Internal Solution Counter**: Because the engine might explore the entire search space looking for extra solutions, it could theoretically finish without currently "holding" a solution. To ensure the solver can still confidently report whether it was successful, it will maintain an internal counter of how many solutions it has found so far. If this counter is greater than zero when the search naturally exhausts, the solver will correctly report success.
4. **Accurate State Retrieval**: When a callback is triggered, developers will need to read the current state of the puzzle. We will add a small safeguard to the high-level `Variable` components ensuring that whenever a variable's value is checked during a callback, it reliably interrogates the underlying engine to fetch the perfectly accurate, real-time solution state.

## Alternatives Considered

During the design phase, several alternative approaches were evaluated and ultimately discarded:

1. **Intelligent Heuristic Scoring:** Instead of finding *all* solutions, we considered injecting word frequency data into the `Trie` dictionary and updating the `Decider` (the logic that chooses which path to explore next). By scoring words, the solver could naturally steer toward the most "human" solution first. We discarded this because it adds significant complexity to the generic solver engine and assumes the user's only goal is word puzzles. The exhaustive search callback is a much more robust, general-purpose feature.
2. **Returning a List of Solutions from `Solve()`:** We discussed modifying the main `Solve()` method to simply return an array (slice) of all valid configurations. We discarded this due to memory constraints. A heavily under-constrained puzzle could have millions of solutions, leading to an immediate out-of-memory crash. The callback mechanism operates like a stream, allowing the developer to evaluate, log, or discard solutions on the fly with a stable memory footprint.
3. **Changing `Solve()` to Return an Integer:** We considered changing the signature of `Solve()` from returning a `bool` (success/failure) to returning an `int` (number of solutions found). We decided against this to maintain backwards compatibility. Existing puzzle logic relies on `if puzzle.Solve(...)` to check for success, and returning a boolean that checks if the internal solution counter is greater than zero seamlessly preserves this behavior.

## Detailed Implementation

To bring this feature to life, the following files will be modified:

### 1. `internal/csp/csp.go`
This is the core constraint solver engine and will require the bulk of the logic changes:
- **`SolutionTracker` Interface:** We will update the signature of `CaptureSolution(p *Problem)` to return a `bool`.
- **`Problem` Struct:** We will add a `solutionsFound int` field to keep track of how many valid endpoints the solver hits.
- **`Problem.recSolve()`:** We will update the logic where the solver detects a success (i.e., when all decisions have been made). It will now increment the `solutionsFound` counter, call the user's `CaptureSolution` callback (if provided), and use the returned boolean to either halt (`return true`) or backtrack by simulating a failure (`return false`).
- **`Problem.Solve()`:** We will change the final return value of this public method to `p.solutionsFound > 0`. This ensures it accurately reports a success even if the developer's callback forced `recSolve` to fully exhaust the search tree and technically "fail" out of the recursive loop.

### 2. `internal/puzzle/puzzle.go`
This is the high-level API where developers interact with typed puzzle values. We need a critical safeguard here for accurate state retrieval during callbacks, without compromising the solver's performance:
- **Callback Interception for State Sync:** Currently, a variable's `currentValues` cache is only synchronized when specific constraints are evaluated. If a developer reads a variable during a solution callback, its cache might be stale. However, updating `Variable.Values()` to dynamically bypass this cache would cause massive map allocations in the hottest path of the constraint solver loop.
Instead, we will intercept the developer's `SolutionTracker` inside `Puzzle.Solve()`. Before passing the callback down, `Puzzle` will inject a wrapper tracker. When a solution is found, this wrapper will manually synchronize the `currentValues` cache for all puzzle variables with the underlying `csp.Problem` state *once*. It will then invoke the developer's original callback. This ensures perfectly accurate state retrieval during the callback while keeping `Variable.Values()` as a cheap `O(1)` read for the constraint engine.

### 3. Testing Changes
To properly test this new behavior and ensure no regressions, the following changes will be made to test files:
- **`cmd/puzzletest/main.go`:** This integration test currently implements `SolutionTracker` (via the `Printer` struct) to print solved Sudoku boards. We will update its `CaptureSolution` method to return a `bool`. To maintain its existing behavior (stopping at the first solution), it will simply `return true`.
- **`internal/csp/csp_test.go` (or `puzzle_test.go`):** We will add a brand new unit test specifically targeting the multiple solution behavior. This test will set up a deliberately under-constrained puzzle (one with a known number of distinct solutions), provide a custom `SolutionTracker` that adds each solution to an array and returns `false`, and assert that exactly the expected number of solutions is captured. We will also assert that the final `Solve()` call correctly returns `true`.
