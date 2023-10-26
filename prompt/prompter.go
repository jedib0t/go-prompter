package prompt

import (
	"context"
	"io"
	"os"
	"time"
)

// Prompter in the interface to create and manage a shell-like interactive
// command prompt.
type Prompter interface {
	// ClearHistory clears all record of previously executed commands.
	ClearHistory()

	// CursorLocation returns the current location of the cursor on the prompt.
	CursorLocation() CursorLocation

	// History returns all the commands executed until now.
	History() []HistoryCommand

	// IsActive returns true if there is an active prompt at the moment.
	IsActive() bool

	// KeyMap returns the current KeyMap for inspection/customization.
	KeyMap() KeyMap

	// NumLines returns the number of lines of text in the current active
	//prompt.
	NumLines() int

	// Prompt prompts. It also watches for cancel KeyEvents on the Context to
	// abort the prompt and return control to client.
	Prompt(ctx context.Context) (string, error)

	// SendInput lets you send strings/runes/KeySequence to the currently active
	// prompt.
	SendInput(a []any, delayBetweenRunes ...time.Duration) error

	// SetAutoCompleter sets up the AutoCompleter that will be used to provide
	// suggestions. Consider this as an auto completer which will be useful for
	// suggesting global stuff like language keywords, or global variables.
	//
	// This should ideally be set ONCE for the lifetime of a Prompter object.
	SetAutoCompleter(global AutoCompleter)

	// SetAutoCompleterContextual sets up the AutoCompleter that will be used to
	// provide suggestions with more priority than the regular AutoCompleter.
	// This is supposed to play the role of providing suggestions like local
	// variables or table names when connected to a database because of a
	// previous command.
	SetAutoCompleterContextual(autoCompleter AutoCompleter)

	// SetCommandShortcuts sets up command shortcuts. For example, if you want
	// to get the prompt input as "/help" when the user presses F1, you'd call
	// this function with the argument:
	//
	//	map[KeySequence]string{
	//	   F1: "/help",
	//	}
	//
	// These shortcuts will take precedence over anything in the prompt and
	// overwrite the contents of the prompt and return control to the caller.
	SetCommandShortcuts(shortcuts map[KeySequence]string)

	// SetDebug enables/disables debug logs/messages in the prompt.
	SetDebug(debug bool)

	// SetFooter sets up the footer line below the prompt line for each render
	// cycle.
	SetFooter(header string)

	// SetFooterGenerator sets up the LineGenerator to be used to generate the
	// footer line below the prompt line for each render cycle.
	SetFooterGenerator(footer LineGenerator)

	// SetHeader sets up the header line above the prompt line for each render
	// cycle.
	SetHeader(header string)

	// SetHeaderGenerator sets up the LineGenerator to be used to generate the
	// header line above the prompt line for each render cycle.
	SetHeaderGenerator(header LineGenerator)

	// SetHistory sets up the history (past inputs) for use when the user wants
	// to move backwards/forwards through the command list
	SetHistory(commands []HistoryCommand)

	// SetHistoryExecPrefix sets up the pattern used to exec command from
	// history. Example (prefix="!"):
	//   - !10 == execute 10th command
	//   - ! == execute last command in history
	SetHistoryExecPrefix(prefix string)

	// SetHistoryListPrefix sets up the prefix used to list commands from
	// history. Example (prefix="!!"):
	//   - !! == list all commands in history;
	//   - !! 10 == list last 10 commands
	SetHistoryListPrefix(prefix string)

	// SetInput sets up the input to be read from the given io.Reader instead of
	// os.Stdin.
	SetInput(input io.Reader)

	// SetKeyMap sets up the KeyMap used for interacting with the user's input.
	SetKeyMap(keyMap KeyMap) error

	// SetOutput sets up the output to go to the given io.Writer instead of
	// os.Stdout.
	SetOutput(output io.Writer)

	// SetPrefix sets up the prefix to use before the prompt.
	//
	// SetPrefix and SetPrefixer override the same property and the last
	// function called takes priority.
	SetPrefix(prefix string)

	// SetPrefixer sets up the prefixer to be called to generate the prefix
	// before the prompt for each render cycle.
	//
	// SetPrefix and SetPrefixer override the same property and the last
	// function called takes priority.
	SetPrefixer(prefixer Prefixer)

	// SetRefreshInterval sets up the minimum interval between consecutive
	// renders. Note that this can be overridden in case of a mandatory override
	// event like a cursor blink.
	SetRefreshInterval(interval time.Duration)

	// SetStyle sets up the Style sheet to be followed for the render.
	SetStyle(s Style)

	// SetSyntaxHighlighter sets up the function that will colorize and
	// highlight keywords in the user-input.
	SetSyntaxHighlighter(highlighter SyntaxHighlighter)

	// SetTerminationChecker sets up the termination checker to check if the
	// user input is done and can be returned to caller on "Terminate" action.
	SetTerminationChecker(checker TerminationChecker)

	// SetWidthEnforcer sets up the function to wrap lines longer than the
	// prompt width.
	SetWidthEnforcer(enforcer WidthEnforcer)

	// Style returns the current Style in use, so it can be modified on the fly
	// in between two prompts.
	Style() *Style
}

// New returns a Prompter than can be used over and over to run a CLI.
//
// It sets some sane defaults:
// - no auto-complete
// - command patterns to invoke history (list old commands, invoke old command)
// - simple prefix "> " for the prompt
// - 60hz refresh rate
// - the default style with a 500ms cursor blink
// - no termination checker (i.e., Enter terminates command)
func New() (Prompter, error) {
	p := &prompt{}
	err := p.SetKeyMap(KeyMapDefault)
	if err != nil {
		return nil, err
	}
	p.SetHistoryExecPrefix(DefaultHistoryExecPrefix)
	p.SetHistoryListPrefix(DefaultHistoryListPrefix)
	p.SetInput(os.Stdin)
	p.SetOutput(os.Stdout)
	p.SetPrefixer(PrefixSimple())
	p.SetRefreshInterval(DefaultRefreshInterval)
	p.SetStyle(StyleDefault)
	p.SetTerminationChecker(TerminationCheckerNone())
	p.SetWidthEnforcer(WidthEnforcerDefault)
	return p, nil
}
