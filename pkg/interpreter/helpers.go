package interpreter

import (
	"bufio"
	"fmt"
	"strings"
)

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

func isImport(input string) bool {
	str := sanitize(input) + " "
	return strings.HasPrefix(str, CommandImport+" ")
}

func isExport(input string) bool {
	str := sanitize(input) + " "
	return strings.HasPrefix(str, CommandExport+" ")
}

func getFileFromCommand(input string) (string, error) {
	input = sanitize(input)
	parts := strings.SplitN(input, " ", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("Missing filename for %s. Syntax: %s [filename]", parts[0], parts[0])
	}
	return parts[1], nil
}

func getHelpString() string {
	return "Syntax:\nFuncDefs: `let [func name] [arg1] [arg2] ... = [expression]\n[expression without vars]\nimport/export [filename]\nOther: help, exit, quit, history, clear"
}

func flush(reader *bufio.Reader) {
	var i int
	for i = 0; i < reader.Buffered(); i++ {
		reader.ReadByte()
	}
}
