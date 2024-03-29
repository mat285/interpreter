package interpreter

const (
	linePrompt  = ">"
	successDone = "Done"

	// CommandQuit is the quit command
	CommandQuit = "quit"
	// CommandExit is the exit command
	CommandExit = "exit"
	// CommandEnv is the env command
	CommandEnv = "env"
	// CommandContext is the context command
	CommandContext = "context"
	// CommandClear is the clear command
	CommandClear = "clear"
	// CommandHistory is the history command
	CommandHistory = "history"
	// CommandHelp is the help command
	CommandHelp = "help"
	// CommandSyntax is the syntax command
	CommandSyntax = "syntax"
	// CommandImport is the import command
	CommandImport = "import"
	// CommandExport is the export command
	CommandExport = "export"
)

var (

	// Commands are all of the commands
	Commands = []string{
		CommandQuit,
		CommandExit,
		CommandEnv,
		CommandContext,
		CommandClear,
		CommandHistory,
		CommandHelp,
		CommandSyntax,
		CommandImport,
		CommandExport,
	}
)
