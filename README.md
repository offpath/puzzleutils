# Puzzleutils

A bunch of tools and utilities to help solve puzzles. Since these are
meant for solving custom puzzles in puzzlehunt settings, the highest
priority is placed on extensibility of code and speed of editing.

The core of these libraries is a constraint satisfaction problem (CSP) solver. To use it effectively, you must understand 4 key concepts:
* Decisions
* Values
* Groups
* Constraints

Decisions & Values
Each puzzle is made up of a set of decisions that must be made, where each individual decision has some number of possible values. A standard sudoku puzzle is essentially 81 decisions, each with 9 possible values. This represents the 81 squares on the grid, each of which can take on a value of 1-9. In a dropquote puzzle, each open box is a decision which must be one of 26 letters. In a nonogram, each box is in the grid is a decision to either be empty or shaded.

Groups & Constraints
Each puzzle type tends to have a natural grouping of decisions to which some constraint is applied to their values. A given decision can exist in multple groupings, and these groupings may have different constraints applied to them. In sudoku, each row, column, and 3x3 square is a group, and all of these groups have the same constraint that their decisions must be unique and cover all values 1-9. In a dropquote puzzle, horizontally connected boxes are groups with the constraint that they must spell a valid word, and boxes in the same colum forum a group with the constraint that they all pull from the same letter bank. In a nonogram, rows and columns have constraints around the number of connected, shaded squares. 

Essentially, if you can model a puzzle using these 4 key concepts, then the CSP solver can use a simple backtracking algorithm to find a solution that meets your constraints if one exists.