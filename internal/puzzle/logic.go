package puzzle

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// Example:
//
// Chef:Freda,Karl,Sonia,Wade
// Dish:cashew tofu,lemon snapper,smoked pork,turkey soup
// Score:42,49,56,63
//
// eq(Dish(Sonia), cashew tofu)
// eq(Score(Karl), plus(Chef(smoked pork), 1))
// eq(Score(Freda, 49))
// eq(Score(Chef(turkey soup)), minus(Score(Sonia), 1 ))
//
// Answers:
// 42, Wade, lemon snapper
// 49, Freda, smoked pork
// 56, Karl, turkey soup
// 63, Sonia, cashew tofu
//
// Supported functions:
// eq
// neq
// or
// gt
// gte
// lt
// lte
// plus
// minus
// <category name>

type category struct {
	name   string
	values []string
}

type LogicPuzzle struct {
	*Puzzle
	categories map[string]*category
	values     map[string]*val
	rules      []astNode
}

type astNode interface {
	Evaluate(lp *LogicPuzzle) valueSet
	TypeCheck() bool
	Type() string
}

type valueSet map[string]bool

type val struct {
	v        string
	category string
	index    int
}

func (v val) String() string {
	return fmt.Sprintf("%s:%s(%d)", v.category, v.v, v.index)
}

func (v val) Evaluate(lp *LogicPuzzle) valueSet {
	return map[string]bool{v.v: true}
}

func (v val) Type() string {
	return v.category
}

func (v val) TypeCheck() bool {
	return true
}

type equality struct {
	left  astNode
	right astNode
	equal bool
}

func (e equality) Evaluate(lp *LogicPuzzle) valueSet {
	// TODO(dneal): returns true if any overlap between value sets
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

func (o or) Evaluate(lp *LogicPuzzle) valueSet {
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

func (c comparison) Evaluate(lp *LogicPuzzle) valueSet {
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

func (p plusMinus) Evaluate(lp *LogicPuzzle) valueSet {
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

func (c connection) Evaluate(lp *LogicPuzzle) valueSet {
	// TODO(dneal)
	return nil
}

func (c connection) TypeCheck() bool {
	return c.arg.Type() != "bool" && c.arg.Type() != "int" && c.arg.Type() != c.Type()
}

func (c connection) Type() string {
	return c.name
}

func (lp *LogicPuzzle) parseExpression(tokens []string, i *int) astNode {
	token := strings.TrimSpace(tokens[*i])
	*i++
	var args []astNode
	if token == "(" {
		*i++
		for {
			args = append(args, lp.parseExpression(tokens, i))
			next := tokens[*i]
			*i++
			if next == ")" {
				break
			}
			if next != "," {
				log.Fatalf("Expected , or )")
			}
		}
	}

	if len(args) == 0 {
		v := lp.values[token]
		if v == nil {
			log.Fatalf("Unknown value: %s", token)
		}
		return v
	}

	if lp.categories[token] != nil {
		if len(args) != 1 {
			log.Fatalf("Expected 1 argument for %s", token)
		}
		return connection{token, args[0]}
	}

	if len(args) != 2 {
		log.Fatalf("Expected 2 arguments for %s", token)
	}

	switch token {
	case "eq":
		fallthrough
	case "neq":
		return equality{args[0], args[1], token == "eq"}
	case "or":
		return or{args[0], args[1]}
	case "gt":
		fallthrough
	case "gte":
		fallthrough
	case "lt":
		fallthrough
	case "lte":
		return comparison{args[0], args[1], token}
	case "plus":
		fallthrough
	case "minus":
		return plusMinus{args[0], args[1], token}
	}

	log.Fatalf("Unknown expression: %s", token)
	return nil
}

func (lp *LogicPuzzle) parseRule(s string) {
	re := regexp.MustCompile("[(,)]|[^(,)]+")
	tokens := re.FindAllString(s, -1)
	i := 0
	astNode := lp.parseExpression(tokens, &i)
	if !astNode.TypeCheck() || astNode.Type() != "bool" {
		log.Fatalf("Invalid rule: %s", s)
	}
	lp.rules = append(lp.rules, astNode)
}

func NewLogicPuzzle(s string) *LogicPuzzle {
	lines := strings.Split(s, "\n")
	result := &LogicPuzzle{}
	i := 0
	for ; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			break
		}
		parts := strings.Split(line, ":")
		result.categories[parts[0]] = &category{parts[0], strings.Split(parts[1], ",")}
	}
	i++
	for ; i < len(lines); i++ {
		result.parseRule(lines[i])
	}
	return result
}
