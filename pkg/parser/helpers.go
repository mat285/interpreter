package parser

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
