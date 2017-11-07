package parser

import (
	"fmt"
	"unicode"
)

// Symbol is a variable
type Symbol string

func parseSymbol(runes []rune, startIdx int) (*Expression, int, error) {
	idx := 0
	for i, r := range runes {
		idx = i
		if unicode.IsLetter(r) {
			continue
		} else if unicode.IsSpace(r) {
			break
		} else if r == '(' || r == '-' {
			break
		} else if isOp(r) {
			break
		} else {
			return nil, -1, invalidSymbolError(r, startIdx+i)
		}
	}
	if idx == len(runes)-1 && unicode.IsLetter(runes[idx]) {
		idx = len(runes)
	}
	e := &Expression{}
	sym := Symbol(string(runes[:idx]))
	if isReserved(string(sym)) {
		return nil, -1, fmt.Errorf("Invalid identifier. `%s` is a reserved word", string(sym))
	}
	e.Symbol = &sym
	return e, idx, nil
}
