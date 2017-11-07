package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/mat285/interpreter/pkg/parser"
)

// Context is the environment context
type context struct {
	funcs   map[string]*parser.Function
	history []string
}

func newContext() *context {
	return &context{
		funcs:   make(map[string]*parser.Function),
		history: make([]string, 0),
	}
}

func (c *context) mapFunc(f *parser.Function) error {
	if f.Name == nil {
		return fmt.Errorf("Cannot map anonymous function")
	}
	c.funcs[*f.Name] = f
	return nil
}

func (c *context) addToHistory(input string) {
	c.history = append(c.history, strings.TrimSpace(input))
}

func (c *context) getHistory() string {
	hs := []string{}
	for i, s := range c.history {
		hs = append(hs, fmt.Sprintf("[%d] %s", i+1, s))
	}
	return strings.Join(hs, "\n")
}

func (c *context) env() string {
	funcs := []string{}
	for _, f := range c.funcs {
		funcs = append(funcs, f.String())
	}
	return strings.Join(funcs, "\n")
}

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
	return str == "quit" || str == "exit"
}

func isEnv(input string) bool {
	str := sanitize(input)
	return str == "env" || str == "context"
}

func isClear(input string) bool {
	str := sanitize(input)
	return str == "clear"
}

func isEval(input string) bool {
	str := sanitize(input)
	return strings.HasPrefix(str, "eval") || strings.HasPrefix(str, "do")
}

func isHistory(input string) bool {
	str := sanitize(input)
	return str == "history"
}

func isHelp(input string) bool {
	str := sanitize(input)
	return str == "help" || str == "syntax"
}

func getHelpString() string {
	return "Syntax:\nFuncDefs: `let [func name] [arg1] [arg2] ... = [expression]\nCalculation: [do | eval] [expression without vars]\nFuncCall [func name]([arg1],[arg2],...)\nOther: help, exit"
}

func isFuncCall(input string) bool {
	str := sanitize(input)
	runes := []rune(str)
	return unicode.IsLetter(runes[0]) && strings.Contains(str, "(")
}

func (c *context) callFunc(call *parser.FunctionCall) (int, error) {
	if f, ok := c.funcs[call.Name]; ok {
		return f.Evaluate(parser.FromFuncMap(c.funcs), call.Inputs...)
	}
	return -1, fmt.Errorf("Unknown function `%s`", call.Name)
}

func (c *context) save(filename string) error {
	return nil
}

func (c *context) load(filename string) error {
	return nil
}

func main() {
	ctx := newContext()

	fmt.Println("Started functional interpreter with new environment. Use quit or exit to end session")
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		ctx.addToHistory(input)
		if isEmpty(input) {
			continue
		} else if isHelp(input) {
			fmt.Println(getHelpString())
			continue
		} else if isQuit(input) {
			os.Exit(0)
		} else if isHistory(input) {
			fmt.Println(ctx.getHistory())
			continue
		} else if isClear(input) {
			ctx.funcs = make(map[string]*parser.Function)
			fmt.Println("Done")
			continue
		} else if isEnv(input) {
			out := ctx.env()
			if len(out) > 0 {
				fmt.Println(ctx.env())
			}
			continue
		} else if isFuncDef(input) {
			f, err := parser.ParseFunction(input, parser.FromFuncMap(ctx.funcs))
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = ctx.mapFunc(f)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("OK", f.String())
		} else {
			input = strings.TrimSpace(input)
			a, err := parser.Parse(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
			val, err := a.EvaluateFull(parser.FromFuncMap(ctx.funcs))
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(val)
		}
	}
}
