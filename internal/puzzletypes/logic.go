package puzzletypes

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/offpath/puzzleutils/internal/constraints"
	"github.com/offpath/puzzleutils/internal/puzzle"
)

type category struct {
	name   string
	values []string
}

type LogicPuzzle struct {
	p             *puzzle.Puzzle
	categories    map[string]*category
	categoryNames []string
	values        map[string]*val
	rules         []astNode
	vars          map[string][]*puzzle.Variable
	numEntities   int
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
	if len(l) == 0 || len(r) == 0 {
		return valueSet{"false": true}
	}
	result := valueSet{}
	for a := range l {
		for b := range r {
			i, j := lp.LookupIndex(a), lp.LookupIndex(b)
			switch c.op {
			case "eq":
				if i == j {
					result["true"] = true
				} else {
					result["false"] = true
				}
			case "neq":
				if i != j {
					result["true"] = true
				} else {
					result["false"] = true
				}
			case "gt":
				if i > j {
					result["true"] = true
				} else {
					result["false"] = true
				}
			case "gte":
				if i >= j {
					result["true"] = true
				} else {
					result["false"] = true
				}
			case "lt":
				if i < j {
					result["true"] = true
				} else {
					result["false"] = true
				}
			case "lte":
				if i <= j {
					result["true"] = true
				} else {
					result["false"] = true
				}
			}
		}
	}
	return result
}

func (c comparison) TypeCheck() bool {
	l, r := c.left.Type(), c.right.Type()
	return l != "int" && l != "bool" && r != "int" && r != "bool"
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
				val := lp.LookupValue(p.left.Type(), i+j)
				if val != "" {
					result[val] = true
				}
			case "minus":
				val := lp.LookupValue(p.left.Type(), i-j)
				if val != "" {
					result[val] = true
				}
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
	result := valueSet{}
	argVals := c.arg.Evaluate(lp)
	sourceCat := c.arg.Type()
	targetCat := c.name

	for s := range argVals {
		for m := 0; m < lp.numEntities; m++ {
			if lp.vars[sourceCat][m].Values()[lp.p.GetStringValue(s)] {
				for t := range lp.vars[targetCat][m].Values() {
					result[t.Str()] = true
				}
			}
		}
	}
	return result
}

func (c connection) TypeCheck() bool {
	return c.arg.Type() != "bool" && c.arg.Type() != "int" && c.arg.Type() != c.Type()
}

func (c connection) Type() string {
	return c.name
}

type LogicRuleConstraint struct {
	lp   *LogicPuzzle
	rule astNode
}

func (c *LogicRuleConstraint) Variables() []*puzzle.Variable {
	var all []*puzzle.Variable
	for _, vars := range c.lp.vars {
		all = append(all, vars...)
	}
	return all
}

func (c *LogicRuleConstraint) Check() bool {
	return c.rule.Evaluate(c.lp)["true"]
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
			// This might be an integer for plus/minus
			if _, err := strconv.Atoi(token); err == nil {
				return val{token, "int", -1}
			}
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
	if index < 0 || index >= len(lp.categories[category].values) {
		return ""
	}
	return lp.categories[category].values[index]
}

func NewLogicPuzzle(p *puzzle.Puzzle, s string) *LogicPuzzle {
	lines := strings.Split(s, "\n")
	result := &LogicPuzzle{
		p:          p,
		categories: map[string]*category{},
		values:     map[string]*val{},
		vars:       map[string][]*puzzle.Variable{},
	}
	i := 0
	for ; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			break
		}
		parts := strings.Split(line, ":")
		catName := parts[0]
		result.categories[catName] = &category{catName, strings.Split(parts[1], ",")}
		result.categoryNames = append(result.categoryNames, catName)
		for j, v := range result.categories[catName].values {
			result.values[v] = &val{v, catName, j}
		}
	}

	result.numEntities = len(result.categories[result.categoryNames[0]].values)

	baseCat := result.categoryNames[0]
	result.vars[baseCat] = make([]*puzzle.Variable, result.numEntities)
	for m := 0; m < result.numEntities; m++ {
		v := p.NewVariable()
		v.SetValueRange(puzzle.ValueSet{p.GetStringValue(result.categories[baseCat].values[m]): true})
		result.vars[baseCat][m] = v
	}

	for k := 1; k < len(result.categoryNames); k++ {
		catName := result.categoryNames[k]
		result.vars[catName] = make([]*puzzle.Variable, result.numEntities)
		allowedVals := puzzle.ValueSet{}
		for _, valStr := range result.categories[catName].values {
			allowedVals[p.GetStringValue(valStr)] = true
		}
		for m := 0; m < result.numEntities; m++ {
			v := p.NewVariable()
			v.SetValueRange(allowedVals)
			result.vars[catName][m] = v
		}
		p.AddConstraint(constraints.NewUniqueConstraint(result.vars[catName], allowedVals))
	}

	i++
	for ; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			result.parseRule(line)
		}
	}

	for _, rule := range result.rules {
		p.AddConstraint(&LogicRuleConstraint{lp: result, rule: rule})
	}

	return result
}

func (lp *LogicPuzzle) String() string {
	var result []string
	for m := 0; m < lp.numEntities; m++ {
		var row []string
		for _, catName := range lp.categoryNames {
			val := lp.vars[catName][m].Values().Value()
			if val != nil {
				row = append(row, val.Str())
			} else {
				row = append(row, "?")
			}
		}
		result = append(result, strings.Join(row, ", "))
	}
	return strings.Join(result, "\n")
}
