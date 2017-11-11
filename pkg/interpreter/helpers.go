package interpreter

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/mat285/interpreter/pkg/parser"
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

type result struct {
	val int
	err error
}

func doAndListen(a *parser.AST, context parser.Context) {
	quit := make(chan int)
	res := make(chan *result)

	go listenForQuit(quit)
	go runEvaluation(a, context, res)
	select {
	case <-quit:
		panic("interrupted")
	case r := <-res:
		if r.err != nil {
			panic(r.err)
		}
		fmt.Println(r.val)
		return
	}
}

func listenForQuit(quit chan int) {
	for {
		// listen for kill signal and push to channel
	}
}

func runEvaluation(a *parser.AST, context parser.Context, res chan *result) {
	defer func() {
		err := recover()
		if err != nil {
			res <- &result{err: fmt.Errorf("%v", err)}
		}
	}()
	val, err := a.EvaluateFull(context)
	res <- &result{val: val, err: err}
}
