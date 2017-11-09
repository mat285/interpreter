package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Expression is an expression
type Expression struct {
	Symbol *Symbol

	Val *Value

	Left  *Expression
	Op    *Operator
	Right *Expression

	Negate bool

	Conditional *Conditional
	Functional  *Functional

	// deferred is for internal use, not for actual expressions
	deferred *DeferredEval
}

// Functional is a function call expression
type Functional struct {
	Name   string
	Inputs []*Expression
}

// String returns a string representation of this expresion
func (exp *Expression) String() string {
	if exp == nil {
		return ""
	}
	if exp.Symbol != nil {
		str := string(*exp.Symbol)
		if exp.Negate {
			return fmt.Sprintf("-%s", str)
		}
		return str
	}
	if exp.Val != nil {
		value := int(*exp.Val)
		if exp.Negate {
			value = (-1) * value
		}
		return strconv.Itoa(value)
	}
	if exp.Conditional != nil {
		return fmt.Sprintf("if %s then %s else %s", exp.Conditional.Predicate.String(), exp.Conditional.True.String(), exp.Conditional.False.String())
	}
	if exp.Functional != nil {
		args := []string{}
		for _, arg := range exp.Functional.Inputs {
			args = append(args, arg.String())
		}
		return fmt.Sprintf("%s(%s)", exp.Functional.Name, strings.Join(args, ","))
	}
	l := exp.Left.String()
	r := exp.Right.String()
	if times().Equal(exp.Op) || divide().Equal(exp.Op) || power().Equal(exp.Op) {
		f := "(%s)"
		l = fmt.Sprintf(f, l)
		r = fmt.Sprintf(f, r)
	}
	str := l + string(*exp.Op) + r
	if plus().Equal(exp.Op) && strings.HasPrefix(r, "-") {
		str = l + r
	}
	if exp.Negate {
		return fmt.Sprintf("-(%s)", str)
	}
	return str
}

// Evaluate evaluates the expression with the given context
func (exp *Expression) Evaluate(context Context) *Expression {
	checkMem()
	defer releaseMem()
	if exp.Val != nil {
		u := *exp.Val
		v := int(u)
		val := Value(v)
		if exp.Negate {
			val = -1 * val
		}
		return &Expression{Val: &val}
	}
	if exp.Symbol != nil {
		if v, ok := context[string(*exp.Symbol)]; ok {
			if v.Value != nil {
				val := Value(*v.Value)
				if exp.Negate {
					val = -1 * val
				}
				return &Expression{Val: &val}
			}
			if v.Symbol != nil {
				sym := Symbol(*v.Symbol)
				return &Expression{Symbol: &sym, Negate: exp.Negate}
			}
		}
		s := string(*exp.Symbol)
		sy := Symbol(s)
		return &Expression{Symbol: &sy, Negate: exp.Negate}
	}

	if exp.Conditional != nil {
		pred := exp.Conditional.Predicate.Evaluate(context)
		if pred.Val != nil {
			v := int(*pred.Val)
			if v != 0 {
				return exp.Conditional.True.Evaluate(context)
			}
			return exp.Conditional.False.Evaluate(context)
		}
	}

	if exp.Functional != nil {
		fn, ok := context[exp.Functional.Name]
		inputs := []*Expression{}
		vals := []ContextVar{}
		for _, arg := range exp.Functional.Inputs {
			input := arg.Evaluate(context)
			inputs = append(inputs, input)
			if input.Val != nil {
				vals = append(vals, FromValue(input.Val))
			}
		}
		if len(vals) == len(inputs) && ok {
			d := &DeferredEval{Function: fn.Function, Inputs: vals}
			// val, err := fn.Evaluate(context, vals...)
			// if err == nil {
			// 	v := Value(val)
			// 	return &Expression{Val: &v}
			// }
			return &Expression{deferred: d}
		}
		fu := &Functional{
			Name:   exp.Functional.Name,
			Inputs: inputs,
		}
		// partial eval
		return &Expression{Functional: fu}
	}

	l := exp.Left.Evaluate(context)
	if l.deferred != nil {
		v, err := l.deferred.Function.Evaluate(context, l.deferred.Inputs...)
		if err == nil {
			val := Value(v)
			l = &Expression{Val: &val}
		}
	}
	r := exp.Right.Evaluate(context)
	if r.deferred != nil {
		v, err := r.deferred.Function.Evaluate(context, r.deferred.Inputs...)
		if err == nil {
			val := Value(v)
			r = &Expression{Val: &val}
		}
	}
	o := *exp.Op
	if l.Val != nil && r.Val != nil {
		lv := *l.Val
		rv := *r.Val
		v := o.Evaluate(lv, rv)
		if exp.Negate {
			v = -1 * v
		}
		return &Expression{Val: &v}
	}

	op := Operator(string(o))
	return &Expression{Left: l, Right: r, Op: &op, Negate: exp.Negate}
}

func buildTable(exp *Expression, table SymbolTable) {
	if exp.Symbol != nil {
		s := string(*exp.Symbol)
		m := make(map[*Symbol]bool)
		if v, ok := table[s]; ok {
			m = v
		}
		m[exp.Symbol] = true
		table[s] = m
	}
	if exp.Left != nil {
		buildTable(exp.Left, table)
	}
	if exp.Right != nil {
		buildTable(exp.Right, table)
	}
}

func parse(runes []rune, startIdx int, onlyFirst bool) (*Expression, int, error) {
	if len(runes) == 0 {
		return nil, -1, fmt.Errorf("No symbols to parse")
	}

	full := &Expression{}
	left := &Expression{}
	right := &Expression{}
	var err error

	// 0: beginning of string
	// 1: have left expression, find op and parse right
	// 2: parsed value
	// 3: parsing function
	state := 0

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch state {
		//beginning of substring
		case 0:
			if unicode.IsSpace(r) {
				continue
			} else if r == 'i' && i+2 < len(runes) && runes[i+1] == 'f' && unicode.IsSpace(runes[i+2]) {
				l, idx, err := parseConditional(runes[i:], i+startIdx)
				if err != nil {
					return nil, 0, err
				}
				left = l
				i = i + idx
				state = 1
			} else if r == '(' {
				// we have a parenthetical expression
				// find the close paren
				close := findParens(runes[i:])
				if close < 0 {
					return nil, -1, fmt.Errorf("Unmatched parenthesis in expression")
				}
				// parse the inner expression and make it the left side
				left, _, err = parse(runes[i+1:i+close], i+1+startIdx, false)
				if err != nil {
					return nil, -1, err
				}
				i = i + close // skip ahead
				state = 1     // we have the left expression, parse the right expression
			} else if unicode.IsLetter(r) {
				// this is a letter, making it a symbol
				exp, idx, err := parseSymbol(runes[i:], i+startIdx)
				if err != nil {
					return nil, -1, err
				}
				left = exp
				i = i + idx - 1
				state = 3
			} else if unicode.IsDigit(r) {
				// This is a value parse it as such
				exp, idx, err := parseValue(runes[i:], i+startIdx)
				if err != nil {
					return nil, -1, err
				}
				left = exp
				i = i + idx - 1
				state = 2

			} else if r == '-' {
				nidx := i + 1
				if nidx >= len(runes) {
					return nil, -1, fmt.Errorf("Reached end of string with incomplete expression")
				}
				next, idx, err := parse(runes[nidx:], nidx+startIdx, true)
				if err != nil {
					return nil, -1, err
				}
				next.Negate = true
				left = next
				left.Op = plus()
				i = i + idx
				state = 1

			} else {
				return nil, -1, invalidSymbolError(r, startIdx+i)
			}

		case 1:
			// if onlyFirst {
			// 	// if only parsing the first symbol, return it
			// 	return left, i, nil
			// }
			// TODO revise this
			if unicode.IsSpace(r) {
				continue
			} else if r == '(' {
				//complicated
				nidx := i
				if nidx >= len(runes) {
					return nil, -1, fmt.Errorf("Reached end of string with incomplete expression")
				}
				next, idx, err := parse(runes[nidx:], nidx+startIdx, true)
				if err != nil {
					return nil, -1, err
				}
				l := left
				left = &Expression{}
				left.Left = l
				left.Right = next
				left.Op = times()
				i = i + idx
				state = 1
			} else if (isPlus(r) || r == '-') && onlyFirst {
				return left, i, nil
			} else if isOp(r) && !isDistributive(r) {
				right, _, err = parse(runes[i+1:], i+1+startIdx, false)
				if err != nil {
					return nil, -1, err
				}
				full.Op = toOp(r)
				full.Left = left
				full.Right = right
				return full, len(runes), nil
			} else if r == '-' {
				nidx := i + 1
				if nidx >= len(runes) {
					return nil, -1, fmt.Errorf("Reached end of string with incomplete expression")
				}
				next, idx, err := parse(runes[nidx:], nidx+startIdx, true)
				if err != nil {
					return nil, -1, err
				}
				next.Negate = true
				l := left
				left = &Expression{}
				left.Left = l
				left.Right = next
				left.Op = plus()
				i = i + idx
				state = 1
			} else if isOp(r) && isDistributive(r) {
				nidx := i + 1
				if nidx >= len(runes) {
					return nil, -1, fmt.Errorf("Reached end of string with incomplete expression")
				}
				next, idx, err := parse(runes[nidx:], nidx+startIdx, true)
				if err != nil {
					return nil, -1, err
				}
				l := left
				left = &Expression{}
				left.Left = l
				left.Right = next
				left.Op = toOp(r)
				i = i + idx
				state = 1
			} else {
				return nil, -1, invalidSymbolError(r, startIdx+i)
			}

		case 2:
			if unicode.IsSpace(r) || isOp(r) || r == '(' || r == '-' {
				// value done
				i = i - 1
				state = 1
			} else if unicode.IsLetter(r) {
				// symbol value exp i.e 2x
				l := &Expression{}
				l.Val = left.Val
				r, idx, err := parseSymbol(runes[i:], i+startIdx)
				if err != nil {
					return nil, -1, err
				}
				left.Left = l
				left.Right = r
				left.Op = times()
				left.Val = nil
				left.Symbol = nil
				i = i + idx - 1
				state = 3
			} else {
				return nil, -1, invalidSymbolError(r, startIdx+i)
			}

		case 3:
			if r == '(' {
				// parse a function
				close := findParens(runes[i:])
				if close < 0 {
					return nil, -1, fmt.Errorf("Unmatched parenthesis in expression")
				}

				f, err := parseFunctionalArgs(runes[i+1:i+close], startIdx+i+1)
				if err != nil {
					return nil, -1, err
				}
				f.Name = string(*left.Symbol)
				left = &Expression{Functional: f}
				i = i + close
				state = 1
			} else {
				state = 1
				i = i - 1
			}

		}
	}
	return left, len(runes) + startIdx, nil
}

func parseFunctionalArgs(runes []rune, startIdx int) (*Functional, error) {
	args := [][]rune{}
	idx := 0
	for i, r := range runes {
		if r == ',' {
			arg := runes[idx:i]
			if len(arg) == 0 {
				return nil, fmt.Errorf("Empty argument given at index %d", startIdx+i)
			}
			args = append(args, arg)
			idx = i + 1
		} else if r == '(' {
			close := findParens(runes[i:])
			i = i + close
		}
	}

	if idx < len(runes) && len(runes[idx:]) > 0 {
		args = append(args, runes[idx:])
	}
	inputs := make([]*Expression, 0, len(args))

	for _, arg := range args {
		exp, _, err := parse(arg, startIdx, false)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, exp)
	}
	return &Functional{
		Inputs: inputs,
	}, nil
}

func invalidSymbolError(r rune, idx int) error {
	return fmt.Errorf("Invalid symbol `%c` at index %d of expression", r, idx)
}
