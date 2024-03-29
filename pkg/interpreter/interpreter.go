package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mat285/interpreter/pkg/parser"
)

// Interpreter is a commandline interpreter
type Interpreter struct {
	Context parser.Context
	History []string
}

// New creates a new interpreter
func New() *Interpreter {
	return &Interpreter{
		Context: parser.NewContext(),
		History: make([]string, 0),
	}
}

// Start starts the interpreter
func (i *Interpreter) Start() {
	fmt.Println("Started functional interpreter with new environment. Use quit or exit to end session. Use help for more information")
	for {
		i.run()
	}
}

func (i *Interpreter) addToHistory(input string) {
	i.History = append(i.History, sanitize(input))
}

func (i *Interpreter) getHistory() string {
	hs := []string{}
	for i, statement := range i.History {
		str := fmt.Sprintf("[%d] %s", i, statement)
		hs = append(hs, str)
	}
	return strings.Join(hs, "\n")
}

func (i *Interpreter) clear() {
	i.Context = parser.NewContext()
}

func (i *Interpreter) env() string {
	vars := []string{}
	for _, v := range i.Context {
		vars = append(vars, v.String())
	}
	return strings.Join(vars, "\n")
}

func (i *Interpreter) mapFunc(f *parser.Function) error {
	if f.Name == nil {
		return fmt.Errorf("Cannot map anonymous function")
	}
	i.Context[*f.Name] = parser.FromFunc(f)
	return nil
}

func (i *Interpreter) importCmd(input string) error {
	file, err := getFileFromCommand(input)
	if err != nil {
		return err
	}
	cmds, err := load(file)
	if err != nil {
		return err
	}
	for _, cmd := range cmds {
		i.interpret(cmd)
	}
	return nil
}

func (i *Interpreter) exportCmd(input string) error {
	file, err := getFileFromCommand(input)
	if err != nil {
		return err
	}
	return save(file, i.Context)
}

func (i *Interpreter) run() {
	defer func() {
		err := recover()
		fmt.Println(err)
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		flush(reader)
		fmt.Print(linePrompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		i.addToHistory(input)
		i.interpret(input)
	}
}

func (i *Interpreter) interpret(input string) {
	if isEmpty(input) {
		return
	} else if isHelp(input) {
		fmt.Println(getHelpString())
	} else if isQuit(input) {
		os.Exit(0)
	} else if isHistory(input) {
		fmt.Println(i.getHistory())
	} else if isClear(input) {
		i.clear()
		fmt.Println(successDone)
	} else if isImport(input) {
		err := i.importCmd(input)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(successDone)
	} else if isExport(input) {
		err := i.exportCmd(input)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(successDone)
	} else if isEnv(input) {
		out := i.env()
		if len(out) > 0 {
			fmt.Println(out)
		}
	} else if isFuncDef(input) {
		f, err := parser.ParseFunction(input, i.Context)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = i.mapFunc(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("OK", f.String())
	} else {
		input = strings.TrimSpace(input)
		a, err := parser.Parse(input)
		if err != nil {
			fmt.Println(err)
			return
		}
		val, err := a.EvaluateFull(i.Context)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(val)
	}
}
