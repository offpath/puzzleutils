package puzzle

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
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
// or
// eq
// neq
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

type or struct {
	left  astNode
	right astNode
}

func (o or) Evaluate(lp *LogicPuzzle) valueSet {
	l, r := o.left.Evaluate(lp), o.right.Evaluate(lp)
	if l["true"] || r["true"] {
		return valueSet{"true": true}
	}
	return valueSet{"false": true}
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
	l, r := c.left.Evaluate(lp), c.right.Evaluate(lp)
	for a := range l {
		for b := range r {
			i, j := lp.LookupIndex(a), lp.LookupIndex(b)
			switch c.op {
			case "eq":
				if i == j {
					return valueSet{"true": true}
				}
			case "neq":
				if i != j {
					return valueSet{"true": true}
				}
			case "gt":
				if i > j {
					return valueSet{"true": true}
				}
			case "gte":
				if i >= j {
					return valueSet{"true": true}
				}
			case "lt":
				if i < j {
					return valueSet{"true": true}
				}
			case "lte":
				if i <= j {
					return valueSet{"true": true}
				}
			}
		}
	}
	return valueSet{"false": true}
}

func (c comparison) TypeCheck() bool {
	l, r := c.left.Type(), c.right.Type()
	return l != "int" && l != "bool" && l == r
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
	l, r := p.left.Evaluate(lp), p.right.Evaluate(lp)
	result := valueSet{}
	for a := range l {
		for b := range r {
			i := lp.LookupIndex(a)
			j, _ := strconv.Atoi(b)
			switch p.op {
			case "plus":
				result[lp.LookupValue(p.left.Type(), i+j)] = true
			case "minus":
				result[lp.LookupValue(p.left.Type(), i-j)] = true
			}
		}
	}
	return result
}

func (p plusMinus) TypeCheck() bool {
	l, r := p.left.Type(), p.right.Type()
	return l != "int" && l != "bool" && r == "int"
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
	if tokens[*i] == "(" {
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
	case "or":
		return or{args[0], args[1]}
	case "eq":
		fallthrough
	case "neq":
		fallthrough
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

func (lp *LogicPuzzle) LookupIndex(v string) int {
	if val := lp.values[v]; val != nil {
		return val.index
	}
	return -1
}

func (lp *LogicPuzzle) LookupValue(category string, index int) string {
	if index < 0 || index > len(lp.categories[category].values) {
		return ""
	}
	return lp.categories[category].values[index]
}

func NewLogicPuzzle(s string) *LogicPuzzle {
	lines := strings.Split(s, "\n")
	result := &LogicPuzzle{
		categories: map[string]*category{},
		values:     map[string]*val{},
	}
	i := 0
	for ; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			break
		}
		parts := strings.Split(line, ":")
		result.categories[parts[0]] = &category{parts[0], strings.Split(parts[1], ",")}
		for i, v := range result.categories[parts[0]].values {
			result.values[v] = &val{v, parts[0], i}
		}
	}
	i++
	for ; i < len(lines); i++ {
		result.parseRule(lines[i])
	}
	return result
}
