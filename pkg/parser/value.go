package parser

import (
	"strconv"
	"unicode"
)

// Value is a value
type Value int

func parseValue(runes []rune, startIdx int) (*Expression, int, error) {
	idx := 0
	for i, r := range runes {
		idx = i
		if unicode.IsDigit(r) {
			continue
		} else if unicode.IsSpace(r) {
			break
		} else if unicode.IsLetter(r) {
			break
		} else if isOp(r) {
			break
		} else if r == '(' || r == '-' {
			break
		} else {
			return nil, -1, invalidSymbolError(r, startIdx+i)
		}
	}
	if idx == len(runes)-1 && unicode.IsDigit(runes[idx]) {
		idx++
	}
	e := &Expression{}
	u, err := strconv.Atoi(string(runes[:idx]))
	if err != nil {
		return nil, -1, err
	}
	val := Value(u)
	e.Val = &val
	return e, idx, nil
}
