package interpreter

import "strings"

func sanitize(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func isFuncDef(input string) bool {
	return strings.HasPrefix(sanitize(input), "let ")
}

func isEmpty(input string) bool {
	return len(sanitize(input)) == 0
}

func isQuit(input string) bool {
	str := sanitize(input)
	return str == CommandQuit || str == CommandExit
}

func isEnv(input string) bool {
	str := sanitize(input)
	return str == CommandEnv || str == CommandContext
}

func isClear(input string) bool {
	str := sanitize(input)
	return str == CommandClear
}

func isHistory(input string) bool {
	str := sanitize(input)
	return str == CommandHistory
}

func isHelp(input string) bool {
	str := sanitize(input)
	return str == CommandHelp || str == CommandSyntax
}

func getHelpString() string {
	return "Syntax:\nFuncDefs: `let [func name] [arg1] [arg2] ... = [expression]\nCalculation: [do | eval] [expression without vars]\nFuncCall [func name]([arg1],[arg2],...)\nOther: help, exit"
}
