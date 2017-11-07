package parser

import (
	"fmt"
	"strings"
	"unicode"
)

// Function is a function
type Function struct {
	Name   *string
	Body   *AST
	Inputs []string
}

// FunctionCall is a function call
type FunctionCall struct {
	Name   string
	Inputs []ContextVar
}

func (f *Function) mapInputs(inputs ...ContextVar) (map[string]ContextVar, error) {
	if len(f.Inputs) != len(inputs) {
		return nil, fmt.Errorf("Input length differs from give inputs. Expected %d inputs, found %d", len(f.Inputs), len(inputs))
	}
	ret := make(map[string]ContextVar)
	for i, input := range inputs {
		ret[f.Inputs[i]] = input
	}
	return ret, nil
}

// Evaluate fully evaluates the function, and errors otherwise
func (f *Function) Evaluate(context map[string]ContextVar, inputs ...ContextVar) (int, error) {
	local, err := f.mapInputs(inputs...)
	if err != nil {
		return -1, err
	}
	return f.Body.EvaluateFull(StitchContext(local, context))
}

// PartialEval partially evaluates the function into another function
func (f *Function) PartialEval(context map[string]ContextVar, inputs ...ContextVar) (*Function, error) {
	local, err := f.mapInputs(inputs...)
	if err != nil {
		return nil, err
	}
	exp := f.Body.Evaluate(StitchContext(local, context))

	return &Function{Body: &AST{Root: exp}, Inputs: f.Inputs[len(inputs):]}, nil
}

// String returns a string representation of this function
func (f *Function) String() string {
	name := ""
	if f.Name != nil {
		name = *f.Name + " = "
	}
	return name + "func(" + strings.Join(f.Inputs, ",") + ") -> " + f.Body.String()
}

// Declaration returns a valid declaration for this function
func (f *Function) Declaration() string {
	if f.Name == nil {
		return ""
	}
	args := strings.Join(f.Inputs, " ")
	body := f.Body.String()
	return fmt.Sprintf("%s %s %s = %s", KeywordLet, *f.Name, args, body)
}

// ParseFunction parses the input string as a function
func ParseFunction(input string, context map[string]ContextVar) (*Function, error) {
	return parseLetFunction(input, context)
}

func parseLetFunction(input string, context map[string]ContextVar) (*Function, error) {
	runes := []rune(input)
	if len(runes) < 5 {
		return nil, fmt.Errorf("Not enough chars for function definition")
	}

	f := &Function{}
	args := map[string]bool{}
	order := []string{}
	state := 0

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch state {
		case 0:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.ToLower(r) == 'l' {
				state = 1
			} else {
				return nil, invalidSymbolError(r, i)
			}
		case 1:
			if unicode.ToLower(r) == 'e' {
				state = 2
			} else {
				return nil, invalidSymbolError(r, i)
			}
		case 2:
			if unicode.ToLower(r) == 't' {
				state = 3
			} else {
				return nil, invalidSymbolError(r, i)
			}
		case 3:
			if unicode.IsSpace(r) {
				state = 4
			} else {
				return nil, invalidSymbolError(r, i)
			}
		case 4:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsLetter(r) {
				name, idx, err := parseFunctionName(runes[i:], i)
				if err != nil {
					return nil, err
				}
				f.Name = &name
				i = i + idx - 1
				state = 5
			} else {
				return nil, invalidSymbolError(r, i)
			}
		case 5:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsLetter(r) {
				symbol, idx, err := parseSymbol(runes[i:], i)
				if err != nil {
					return nil, err
				}
				arg := string(*symbol.Symbol)
				if _, ok := args[arg]; ok {
					return nil, fmt.Errorf("Duplicate input `%s`", arg)
				}
				args[arg] = true
				order = append(order, arg)
				i = i + idx - 1
			} else if r == '=' {
				body, err := Parse(string(runes[i+1:]))
				if err != nil {
					return nil, err
				}
				f.Body = body
				i = len(runes)
				break
			} else {
				return nil, invalidSymbolError(r, i)
			}
		}
	}

	if f.Body == nil {
		return nil, fmt.Errorf("Missing function body")
	}

	f.Inputs = order
	return f, f.validate()
}

func parseFunctionName(runes []rune, startIdx int) (string, int, error) {
	idx := 0
	for i, r := range runes {
		idx = i
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			continue
		} else if unicode.IsSpace(r) {
			break
		} else {
			return "", -1, invalidSymbolError(r, i+startIdx)
		}
	}
	return string(runes[:idx]), idx, nil
}

func (f *Function) validate() error {
	if f.Name != nil {
		name := strings.ToLower(*f.Name)
		for _, word := range Keywords {
			if word == name {
				return fmt.Errorf("Invalid function name. Name matches reserved word `%s`", word)
			}
		}
	}

	symbols := f.Body.Symbols()

	args := map[string]bool{}

	for _, in := range f.Inputs {
		args[in] = true
	}

	for _, sym := range symbols {
		if !args[sym] {
			return fmt.Errorf("Unknown symbol `%s` is not defined", sym)
		}
	}
	return nil
}
