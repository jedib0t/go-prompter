package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/jedib0t/go-prompter/powerline"
	"github.com/jedib0t/go-prompter/prompt"
	"github.com/muesli/termenv"
)

var (
	flagDebug = flag.Bool("debug", false, "Enable Debug logging?")

	commandShortcuts = map[prompt.KeySequence]string{
		prompt.CtrlC:  "quit",
		prompt.Escape: "quit",
	}

	segmentHost   = powerline.Segment{}
	segmentUser   = powerline.Segment{}
	segmentCmdNum = powerline.Segment{}
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
	p.SetPrefixer(generatePowerlinePrefixer())

	fmt.Println("Simple Prompt: (ctrl+c  to quit)")
	cmdNum := 0
	for {
		cmdNum++
		segmentCmdNum.SetContent(fmt.Sprintf("%d", cmdNum))

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

func generatePowerlinePrefixer() prompt.Prefixer {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	userObj, err := user.Current()
	if err != nil {
		userObj = &user.User{Username: "username"}
	}

	// The prefix is basically a Powerline made up of multiple segments; the
	// following sets up all the segments we are going to use.
	segmentHost.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(16), Background: termenv.ANSI256Color(220)})
	segmentHost.SetContent(hostname)
	segmentUser.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(17)})
	segmentUser.SetContent(userObj.Username)
	segmentCmdNum.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(16), Background: termenv.ANSI256Color(147)})

	p := powerline.Powerline{}
	p.Append(&segmentHost)
	p.Append(&segmentUser)
	p.Append(&segmentCmdNum)
	p.SetStyle(powerline.StylePatched)
	return func() string {
		return p.Render(0) + " "
	}
}
