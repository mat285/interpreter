package parser

import (
	"fmt"
)

// AST is an abstract syntax tree
type AST struct {
	Root        *Expression
	SymbolTable SymbolTable
}

// SymbolTable is a table of symbols
type SymbolTable map[string]map[*Symbol]bool

// String returns a string of the ast
func (a *AST) String() string {
	return a.Root.String()
}

// Symbols returns the list of symbols in this tree
func (a *AST) Symbols() []string {
	table := SymbolTable(make(map[string]map[*Symbol]bool))
	buildTable(a.Root, table)
	a.SymbolTable = table
	syms := make([]string, 0, len(a.SymbolTable))
	for k := range a.SymbolTable {
		syms = append(syms, k)
	}
	return syms
}

// eval evaluates the tree with the given symbol value map
func (a *AST) eval(table ...Context) *Expression {
	t := make(Context)
	if len(table) > 0 && table[0] != nil {
		t = table[0]
	}
	return a.Root.Evaluate(t)
}

// Evaluate evaluates the tree with the given context
func (a *AST) Evaluate(table ...Context) *Expression {
	t := make(Context)
	if len(table) > 0 && table[0] != nil {
		t = table[0]
	}
	exp := a.Root.Evaluate(t)
	for exp.deferred != nil {
		exp = exp.deferred.Eval(t)
	}
	return exp
}

// EvaluateFull evaluates this expression down to in an int if possible, or fails
func (a *AST) EvaluateFull(table ...Context) (int, error) {
	exp := a.Evaluate(table...)
	if exp == nil || exp.Val == nil {
		return -1, fmt.Errorf("Could not fully evaluate the expression. Variables still remain: %s", exp.String())
	}
	return int(*exp.Val), nil
}

// evaluateFull evaluates this expression down to in an int if possible, or fails
func (a *AST) evalFull(table ...Context) (int, error) {
	exp := a.eval(table...)
	if exp == nil || exp.Val == nil {
		return -1, fmt.Errorf("Could not fully evaluate the expression. Variables still remain: %s", exp.String())
	}
	return int(*exp.Val), nil
}

// Parse parses the expression into an abstract syntax tree
func Parse(exp string) (*AST, error) {
	e, _, err := parse([]rune(exp), 0, false)
	if err != nil {
		return nil, err
	}
	table := SymbolTable(make(map[string]map[*Symbol]bool))
	buildTable(e, table)
	return &AST{
		Root:        e,
		SymbolTable: table,
	}, nil
}
