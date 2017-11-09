package parser

var stackCalls = 0

func isReserved(s string) bool {
	for _, k := range Keywords {
		if s == k {
			return true
		}
	}
	return false
}

// find the index with the close parens
func findParens(runes []rune) int {
	depth := 0
	for i, r := range runes {
		if r == '(' {
			depth++
		} else if r == ')' {
			depth--
		}

		if depth == 0 {
			return i
		}
	}
	return -1
}

func checkMem() {
	stackCalls++
	if stackCalls >= maxStackCalls {
		panic("stack overflow: too many recursive calls")
	}
}

func releaseMem() {
	stackCalls--
}
