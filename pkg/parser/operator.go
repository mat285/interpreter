package parser

import (
	"math"
)

// Operator is an operator
type Operator string

const (
	// Plus is addition
	Plus Operator = "+"
	// Times is multiplication
	Times Operator = "*"
	// Divided is division
	Divided Operator = "/"
	// Power is a exponentiation
	Power Operator = "^"
	// GreaterThan is greater than
	GreaterThan Operator = ">"
	// LessThan is less than
	LessThan Operator = "<"
	// Or is or
	Or Operator = "|"
	// And is and
	And Operator = "&"
	// Equal is equality
	Equal Operator = "="
)

// Equal tests for equality
func (o *Operator) Equal(op *Operator) bool {
	return o == op || string(*o) == string(*op)
}

// Distributive tests for distribution
func (o *Operator) Distributive() bool {
	return o.Equal(times()) || o.Equal(divide()) || o.Equal(power()) || o.Equal(and())
}

// Evaluate evaluates this operator
func (o Operator) Evaluate(v1, v2 Value) Value {
	i1 := int(v1)
	i2 := int(v2)
	switch o {
	case Plus:
		return Value(i1 + i2)
	case Times:
		return Value(i1 * i2)
	case Divided:
		return Value(i1 / i2)
	case Power:
		return Value(int(math.Pow(float64(i1), float64(i2))))
	case GreaterThan:
		if v1 > v2 {
			return Value(1)
		}
		return Value(0)
	case LessThan:
		if v1 < v2 {
			return Value(1)
		}
		return Value(0)
	case Or:
		if v1 != Value(0) || v2 != Value(0) {
			return Value(1)
		}
		return Value(0)
	case And:
		if v1 != Value(0) && v2 != Value(0) {
			return Value(1)
		}
		return Value(0)
	case Equal:
		if v1 == v2 {
			return Value(1)
		}
		return Value(0)
	}
	panic("Unknown op")
}

func plus() *Operator {
	op := Plus
	return &op
}

func times() *Operator {
	op := Times
	return &op
}

func divide() *Operator {
	op := Divided
	return &op
}

func power() *Operator {
	op := Power
	return &op
}

func gt() *Operator {
	op := GreaterThan
	return &op
}

func lt() *Operator {
	op := LessThan
	return &op
}

func or() *Operator {
	op := Or
	return &op
}

func and() *Operator {
	op := And
	return &op
}

func eq() *Operator {
	op := Equal
	return &op
}

func isOp(r rune) bool {
	return isTimes(r) || isPlus(r) || isDivided(r) || isPower(r) || isGT(r) || isLT(r) || isAnd(r) || isOr(r) || isEq(r)
}

func isDistributive(r rune) bool {
	return toOp(r).Distributive()
}

func toOp(r rune) *Operator {
	o := Operator(string(r))
	return &o
}

func isTimes(r rune) bool {
	return string(r) == string(Times)
}

func isPlus(r rune) bool {
	return string(r) == string(Plus)
}

func isDivided(r rune) bool {
	return string(r) == string(Divided)
}

func isPower(r rune) bool {
	return string(r) == string(Power)
}

func isGT(r rune) bool {
	return string(r) == string(GreaterThan)
}

func isLT(r rune) bool {
	return string(r) == string(LessThan)
}

func isOr(r rune) bool {
	return string(r) == string(Or)
}

func isAnd(r rune) bool {
	return string(r) == string(And)
}

func isEq(r rune) bool {
	return string(r) == string(Equal)
}
