package parser

import (
	"fmt"
)

func parseIfStatement(runes []rune, startIdx int) (*Expression, int, error) {
	lenIf := len([]rune(KeywordIf))
	lenThen := len([]rune(KeywordThen))
	lenElse := len([]rune(KeywordElse))

	close := findIfThenClose(runes) - 4
	if close <= 0 {
		return nil, -1, fmt.Errorf("Mismatched if-then statement")
	}

	pred, _, err := parse(runes[lenIf:close-1], startIdx, false)
	if err != nil {
		return nil, -1, err
	}

	next := findThenElseClose(runes[close:]) + close - lenElse
	if next <= close {
		return nil, -1, fmt.Errorf("Mismatched then-else statement")
	}

	then, _, err := parse(runes[close+lenThen:next-1], startIdx+close, false)
	if err != nil {
		return nil, -1, err
	}

	l := next + lenElse
	if l >= len(runes) {
		return nil, -1, fmt.Errorf("Missing else expression")
	}

	last, _, err := parse(runes[l:], startIdx+l, false)
	if err != nil {
		return nil, -1, err
	}
	cond := &Conditional{
		Predicate: pred,
		True:      then,
		False:     last,
	}
	return &Expression{Conditional: cond}, len(runes), nil
}

func findIfThenClose(runes []rune) int {
	depth := 0
	state := 0
	idx := 0
	for i := 0; i < len(runes); i++ {
		idx = i
		r := runes[i]
		switch state {
		case 0:
			if r == 'i' {
				state = 3
			} else {
				return -1
			}
		case 1:
			if r == ' ' {
				state = 2
			}
		case 2:
			if r == 'i' {
				state = 3
			} else if r == 't' {
				state = 5
			}
		case 3:
			if r == 'f' {
				state = 4
			} else {
				state = 1
			}
		case 4:
			if r == ' ' {
				depth++
			}
			state = 2

		case 5:
			if r == 'h' {
				state = 6
			} else {
				state = 1
			}
		case 6:
			if r == 'e' {
				state = 7
			} else {
				state = 1
			}
		case 7:
			if r == 'n' {
				state = 8
			} else {
				state = 1
			}
		case 8:
			if r == ' ' {
				depth--
			}
			state = 1
		}
		if i > 3 && depth == 0 {
			break
		}
	}
	if depth != 0 {
		return -1
	}
	return idx
}

func findThenElseClose(runes []rune) int {
	depth := 0
	state := 0
	idx := 0
	for i := 0; i < len(runes); i++ {
		idx = i
		r := runes[i]
		switch state {
		case 0:
			if r == 't' {
				state = 3
			} else {
				return -1
			}
		case 1:
			if r == ' ' {
				state = 2
			}
		case 2:
			if r == 't' {
				state = 3
			} else if r == 'e' {
				state = 7
			}
		case 3:
			if r == 'h' {
				state = 4
			} else {
				state = 1
			}
		case 4:
			if r == 'e' {
				state = 5
			} else {
				state = 1
			}

		case 5:
			if r == 'n' {
				state = 6
			} else {
				state = 1
			}
		case 6:
			if r == ' ' {
				depth++
			}
			state = 2
		case 7:
			if r == 'l' {
				state = 8
			} else {
				state = 1
			}
		case 8:
			if r == 's' {
				state = 9
			} else {
				state = 1
			}
		case 9:
			if r == 'e' {
				state = 10
			} else {
				state = 1
			}
		case 10:
			if r == ' ' {
				depth--
			}
			state = 2
		}
		if i > 5 && depth == 0 {
			break
		}
	}
	if depth != 0 {
		return -1
	}
	return idx
}
