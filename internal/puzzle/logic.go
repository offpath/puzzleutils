package puzzle

import (
	"fmt"
	"strings"
)

// Example:
//
// Chef:Freda,Karl,Sonia,Wade
// Dish:cashew tofu,lemon snapper,smoked pork,turkey soup
// Score:42,49,56,63
//
// eq(Dish(Sonia), cashew tofu)
// eq(Score(Karl), plus(Chef(smoked pork), 7))
// eq(Score(Freda, 49))
// eq(Score(Chef(turkey soup)), minus(Score(Sonia), 7))
//
// Answers:
// 42, Wade, lemon snapper
// 49, Freda, smoked pork
// 56, Karl, turkey soup
// 63, Sonia, cashew tofu

type category struct {
	name   string
	values []string
}

type LogicPuzzle struct {
	*Puzzle
	categories map[string]*category
}

type astNode interface {
	Evaluate() valueSet
	TypeCheck() bool
	Type() string
}

type valueSet map[string]bool

type val struct {
	v string
	t string
}

func (v val) String() string {
	return fmt.Sprintf("%s:%s", v.t, v.v)
}

func (v val) Evaluate() valueSet {
	return map[string]bool{v.v: true}
}

func (v val) Type() string {
	return v.t
}

func (v val) TypeCheck() bool {
	return true
}

// TODO(dneal): Functions to support:
// eq
// neq
// or
// gt
// gte
// lt
// lte
// <category name>
// plus
// minus

type equality struct {
	left  astNode
	right astNode
	equal bool
}

func (e equality) Evalate() valueSet {
	// TODO(dneal)
	return nil
}

func (e equality) TypeCheck() bool {
	return e.left.Type() == e.right.Type()
}

func (e equality) Type() string {
	return "bool"
}

type or struct {
	left  astNode
	right astNode
}

func (o or) Evalate() valueSet {
	// TODO(dneal)
	return nil
}

func (o or) TypeCheck() bool {
	return o.left.Type() == "bool" && o.right.Type() == "bool"
}

func (o or) Type() string {
	return "bool"
}

type comparison struct {
	left  astNode
	right astNode
	op    string
}

func (c comparison) Evalate() valueSet {
	// TODO(dneal)
	return nil
}

func (c comparison) TypeCheck() bool {
	return c.left.Type() == c.right.Type()
}

func (c comparison) Type() string {
	return "bool"
}

type plusMinus struct {
	left  astNode
	right astNode
	op    string
}

func (p plusMinus) Evalate() valueSet {
	// TODO(dneal)
	return nil
}

func (p plusMinus) TypeCheck() bool {
	return p.left.Type() == p.right.Type()
}

func (p plusMinus) Type() string {
	return p.left.Type()
}

type connection struct {
	name string
	arg  astNode
}

func (c connection) Evalate() valueSet {
	// TODO(dneal)
	return nil
}

func (c connection) TypeCheck() bool {
	return c.arg.Type() != "bool" && c.arg.Type() != "int" && c.arg.Type() != c.Type()
}

func (c connection) Type() string {
	return c.name
}

func NewLogicPuzzle(s string) *LogicPuzzle {
	lines := strings.Split(s, "\n")
	result := &LogicPuzzle{}
	for _, line := range lines {
		if line == "" {
			break
		}
		parts := strings.Split(line, ":")
		result.categories[parts[0]] = &category{parts[0], strings.Split(parts[1], ",")}
	}
	return result
}
