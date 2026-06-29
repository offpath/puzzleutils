# Tracker

This directory contains utility implementations for monitoring the progress of the `csp` solver during execution. Because solving complex constraint satisfaction problems can sometimes take an unpredictable amount of time (due to the exponential nature of backtracking search), these trackers allow developers and users to see how many decisions the solver is actively making in real-time.

## Files

### `tracker.go`
This is the only file in the directory. It provides concrete implementations of the `csp.DecisionTracker` interface, which acts as a callback fired every time the solver makes a guess/decision. It exposes two main functions:

- **`PrintEveryN(n int)`**: Creates a tracker that prints the total number of decisions made to standard output exactly every `n` decisions.
- **`PrintEveryLogN(n int)`**: Creates a tracker that prints logarithmically. For example, if `n` is 10, it will print at 1, 2, ..., 10, then 20, ..., 100, then 200, ..., 1000. This is highly useful for monitoring progress over both fast and extremely long-running solves without spamming the console with output.
