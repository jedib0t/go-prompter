package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-prompter/prompt"
	"github.com/muesli/termenv"
)

var (
	flagDebug       = flag.Bool("debug", false, "Enable Debug logging?")
	flagDisableTime = flag.Bool("disable-time", false, "Disable Timestamp?")

	colorPromptSymbol = prompt.Color{
		Foreground: termenv.ANSI256Color(154),
		Background: termenv.BackgroundColor(),
	}
	colorPromptTime = prompt.Color{
		Foreground: termenv.ANSI256Color(8),
		Background: termenv.BackgroundColor(),
	}

	commandShortcuts = map[prompt.KeySequence]string{
		prompt.CtrlC:  "quit",
		prompt.Escape: "quit",
	}
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p, err := prompt.New()
	if err != nil {
		fmt.Printf("ERROR: : %v", err)
		os.Exit(1)
	}
	p.SetCommandShortcuts(commandShortcuts)
	p.SetDebug(*flagDebug)
	p.SetPrefixer(prefixer)

	fmt.Println("Simple Prompt: (ctrl+c  to quit)")
	for {
		input, err := p.Prompt(ctx)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			break
		}
		switch strings.ToLower(input) {
		case "bye", "exit", "quit":
			fmt.Printf("Bye!\n")
			os.Exit(0)
		default:
			fmt.Printf("> Executing: %#v\n", input)
		}
		fmt.Println()
	}
}

func prefixer() string {
	if *flagDisableTime {
		return fmt.Sprintf("%s ", colorPromptSymbol.Sprint("$"))
	}

	return fmt.Sprintf("%s %s ",
		colorPromptTime.Sprintf("[%s]", time.Now().Format(time.TimeOnly)),
		colorPromptSymbol.Sprint("$"),
	)
}
