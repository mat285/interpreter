package parser

const (
	// KeywordIf is the keyword for if
	KeywordIf = "if"
	// KeywordThen is the keyword for then
	KeywordThen = "then"
	// KeywordElse is the keyword for else
	KeywordElse = "else"
	// KeywordLet is the keyword for let
	KeywordLet = "let"

	maxStackCalls = 500000
)

var (
	// Keywords are the reserved words for expressions
	Keywords = []string{KeywordIf, KeywordElse, KeywordThen, KeywordLet}
)
