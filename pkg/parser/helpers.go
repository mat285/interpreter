package parser

func isReserved(s string) bool {
	for _, k := range Keywords {
		if s == k {
			return true
		}
	}
	return false
}
