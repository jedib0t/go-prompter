package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-prompter/prompt"
)

var (
	flagDebug          = flag.Bool("debug", false, "Enable Debug logging?")
	flagDemo           = flag.Bool("demo", false, "Auto-execute a few commands in the prompt?")
	flagDisableLineNum = flag.Bool("disable-line-num", false, "Disable Line numbers?")
	flagHeightMax      = flag.Uint("height-max", 5, "Maximum Height (excluding title); 0==no-limit")
	flagHeightMin      = flag.Uint("height-min", 1, "Minimum Height (excluding title); 0==no-limit")
	flagWidthMax       = flag.Uint("width-max", 0, "Maximum Terminal Width to use (0 for full length)")
	flagStyle          = flag.String("style", "monokai", "Chroma Style to use for syntax highlighting")
	flagTimeoutSecs    = flag.Uint("timeout", 300, "Number of seconds to timeout after.")

	shortcuts = map[prompt.KeySequence]string{
		prompt.F1:     "/help",
		prompt.Escape: "/quit",
		prompt.CtrlC:  "/quit",
	}
)

func main() {
	flag.Parse()

	// You can use a context.WithCancel to use when you want to abort the prompt
	// midway through user input for whatever reason. Here, we demo this using
	// a timeout on the prompt.
	timeout := time.Second * time.Duration(*flagTimeoutSecs)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// KeyMap can be customized for supported actions.
	keyMap := prompt.KeyMapMultiLine
	keyMap.Insert.Abort = append(keyMap.Insert.Abort, prompt.F10)

	// Syntax highlight the input using the Chroma library. A shorthand for this
	// for SQL would be "prompt.SyntaxHighlighterSQL" - it uses some preselected
	// defaults on "prompt.SyntaxHighlighterChroma".
	syntaxHighlighter, err := prompt.SyntaxHighlighterChroma("sql", "terminal256", *flagStyle)
	if err != nil {
		fmt.Printf("ERROR: failed to initialize syntax highlighter: %v", err)
		os.Exit(1)
	}

	// Get a new prompter, and for demo reasons we are not going to be using
	// prompt.SQL() which does a few of the following "Set"s by default.
	p, err := prompt.New()
	if err != nil {
		fmt.Printf("ERROR: failed to initialize prompt: %v", err)
		os.Exit(1)
	}
	p.SetAutoCompleter(prompt.AutoCompleteSQLKeywords())
	p.SetAutoCompleterContextual(prompt.AutoCompleteSimple(tableAndColumnNames, true))
	p.SetCommandShortcuts(shortcuts)
	p.SetDebug(*flagDebug)
	if !*flagDemo {
		p.SetHistory(history)
	}
	p.SetHistoryExecPrefix("!")
	p.SetHistoryListPrefix("!!")
	err = p.SetKeyMap(keyMap)
	if err != nil {
		fmt.Printf("ERROR: failed to initialize key-map: %v", err)
		os.Exit(1)
	}
	p.SetPrefixer(prompt.PrefixNone())
	p.SetSyntaxHighlighter(syntaxHighlighter)
	p.SetTerminationChecker(prompt.TerminationCheckerSQL())
	p.Style().Dimensions.HeightMax = *flagHeightMax
	p.Style().Dimensions.HeightMin = *flagHeightMin
	p.Style().Dimensions.WidthMax = *flagWidthMax
	if !*flagDisableLineNum {
		p.Style().LineNumbers = prompt.StyleLineNumbersEnabled
	}

	// We are going to use the "powerline" package in this library to generate a
	// fancy header bar above the prompt, loaded with things like the current
	// command number, cursor location, timestamp, etc.
	p.SetHeaderGenerator(generatePowerline())
	// Asynchronously update the cursor location and timestamp continuously.
	go func() {
		tick := time.Tick(time.Second / 5)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick:
				segmentCursor.SetContent(p.CursorLocation().String())
				segmentTime.SetContent(time.Now().Format(time.TimeOnly))
			}
		}
	}()

	// turn on demo mode if asked for
	if *flagDemo {
		go runDemo(p)
	}

	// Prompt the user and handle each input in a loop until we are done for any
	// reason (user wants to quit, etc.).
	for {
		// Update the # of the command we are handling on the title bar
		segmentCmdNum.SetContent(fmt.Sprint(len(p.History()) + 1))

		// Prompt
		input, err := p.Prompt(ctx)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err.Error())
			os.Exit(1)
		}

		// Handle the input
		cmd, _ := getCommandAndArgs(input)
		switch strings.ToLower(cmd) {
		case "/?", "/help":
			printHelp()
		case "/clear":
			p.ClearHistory()
			fmt.Println("Cleared history.")
		case "/quit":
			fmt.Println("Bye!")
			os.Exit(0)
		default:
			// pretend we talk to a real database and output real data
			printDummyOutput(input)
		}
		fmt.Println()
	}
}

var (
	reCommand        = regexp.MustCompile(`^\s*(\S+)\s*$`)
	reCommandAndArgs = regexp.MustCompile(`^\s*(\S+)\s*(.*)$`)
	reSqlComments    = regexp.MustCompile(`(/\*.*\*/|--[^\n]*\n|--[^\n]*$)`)
)

func getCommandAndArgs(input string) (string, []string) {
	input = reSqlComments.ReplaceAllString(input, "")
	input = strings.TrimSpace(input)

	if matches := reCommand.FindStringSubmatch(input); len(matches) == 2 {
		return matches[1], nil
	}
	if matches := reCommandAndArgs.FindStringSubmatch(input); len(matches) == 3 {
		return matches[1], strings.Split(matches[2], " ")
	}
	return input, nil
}

func printDummyOutput(input string) {
	switch input {
	case history[0].Command: // "select * from employees where id = 1;":
		tw := tableWriter()
		tw.AppendHeader(table.Row{"ID", "First Name", "Last Name", "Salary", "Notes"})
		tw.AppendRow(table.Row{1, "Night", "King", 10000, "Has horns!"})
		tw.SetCaption("Returned 1 row in 0.001s.")
		fmt.Println(tw.Render())
	case history[1].Command: // "insert into employees (first_name, last_name, salary) values\n  ('Arya', 'Stark', 3000),\n  ('Jon', 'Snow', 2000),\n  ('Tyrion', 'Lannister', 5000);"
		fmt.Println("Inserted 3 records in 0.015s.")
	case history[2].Command: // "select * from employees where salary between 1000 and 6000 order by id;"
		tw := tableWriter()
		tw.AppendHeader(table.Row{"ID", "First Name", "Last Name", "Salary", "Notes"})
		tw.AppendRow(table.Row{2, "Arya", "Start", 3000, "Not today."})
		tw.AppendRow(table.Row{3, "Jon", "Snow", 2000, "Knows nothing."})
		tw.AppendRow(table.Row{4, "Tyrion", "Lannister", 5000, "Pays his debts."})
		tw.SetCaption("Returned 3 rows in 0.003s.")
		fmt.Println(tw.Render())
	case history[3].Command: // "delete from employees where salary between 1000 and 6000;"
		fmt.Println("Deleted 3 records in 0.013s.")
	default:
		fmt.Printf("> Pretending to execute: %#v\n", input)
	}
}

func printHelp() {
	fmt.Println(`SQL Prompt demo using github.com/jedib0t/go-prompter.

* /?, /help    Prints this help text.
* /clear       Clears History.
* /quit        Exits the prompt.`)
}

func tableWriter() table.Writer {
	tw := table.NewWriter()
	tw.SetStyle(table.StyleLight)
	return tw
}
