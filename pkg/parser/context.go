package parser

import "fmt"

// Context is the context under which things are parsed
type Context map[string]ContextVar

// ContextVar is a type that can be a symbol or a value
type ContextVar struct {
	*Symbol
	*Value
	*Function
}

// NewContext creates a new context
func NewContext() Context {
	return make(Context)
}

// Clone clones the context
func (c Context) Clone() Context {
	ret := make(Context)
	for k, v := range c {
		ret[k] = v
	}
	return ret
}

// Source returns a source code string for the functions in this context
func (c Context) Source() string {
	ret := ""
	for _, v := range c {
		if v.Function != nil {
			ret = fmt.Sprintf("%s\n", v.Function.Declaration())
		}
	}
	return ret
}

func (c ContextVar) String() string {
	if c.Symbol != nil {
		return string(*c.Symbol)
	}
	if c.Value != nil {
		return fmt.Sprintf("%d", *c.Value)
	}
	if c.Function != nil {
		return c.Function.String()
	}
	return "" // make clearer
}

// ToExpression converts a context var to an expression
func ToExpression(v ContextVar) *Expression {
	if v.Symbol != nil {
		return &Expression{Symbol: v.Symbol}
	}
	if v.Value != nil {
		return &Expression{Val: v.Value}
	}
	syms := []*Expression{}
	for _, s := range v.Function.Inputs {
		sym := Symbol(s)
		syms = append(syms, &Expression{Symbol: &sym})
	}
	fn := &Functional{
		Name:   *v.Function.Name,
		Inputs: syms,
	}
	return &Expression{Functional: fn}
}

// FromSymbol returns a Symbol ContextVar
func FromSymbol(s *Symbol) ContextVar {
	return ContextVar{s, nil, nil}
}

// FromValue returns a Value ContextVar
func FromValue(v *Value) ContextVar {
	return ContextVar{nil, v, nil}
}

// FromFunc Creates a context var from a function
func FromFunc(f *Function) ContextVar {
	return ContextVar{nil, nil, f}
}

// FromIntMap transforms the int map into t symbolvalue map
func FromIntMap(input map[string]int) Context {
	ret := make(Context)
	for k, v := range input {
		val := Value(v)
		ret[k] = FromValue(&val)
	}
	return ret
}

// FromFuncMap returns a context map from func map
func FromFuncMap(input map[string]*Function) Context {
	ret := make(Context)
	for k, v := range input {
		ret[k] = FromFunc(v)
	}
	return ret
}

// StitchContext stitches the local and global context with local over global
func StitchContext(local Context, global Context) Context {
	ret := make(Context)
	for k, v := range local {
		ret[k] = v
	}
	for k, v := range global {
		if _, ok := ret[k]; !ok {
			ret[k] = v
		}
	}
	return ret
}
