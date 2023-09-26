package main

import (
	"os"
	"os/user"
	"time"

	"github.com/jedib0t/go-prompter/powerline"
	"github.com/jedib0t/go-prompter/prompt"
	"github.com/muesli/termenv"
)

var (
	segmentHost   = powerline.Segment{}
	segmentUser   = powerline.Segment{}
	segmentDBName = powerline.Segment{}
	segmentDBType = powerline.Segment{}
	segmentCmdNum = powerline.Segment{}
	segmentCursor = powerline.Segment{}
	segmentTime   = powerline.Segment{}
)

func generatePowerline() prompt.LineGenerator {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	userObj, err := user.Current()
	if err != nil {
		userObj = &user.User{Username: "username"}
	}

	// The title bar is basically a Powerline made up of multiple segments; the
	// following sets up all the segments we are going to use.
	segmentHost.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(16), Background: termenv.ANSI256Color(220)})
	segmentHost.SetContent(hostname)
	segmentUser.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(17)})
	segmentUser.SetContent(userObj.Username)
	// For the DB segments, we want it to be colored differently based on which
	// DB we are connected to. Don't set a color and let Powerline generate a
	// color based on the content value.
	segmentDBName.SetContent("foo.user@bar.db")
	segmentDBType.SetContent("MySQL")
	segmentCmdNum.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(16), Background: termenv.ANSI256Color(147)})
	segmentTime.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(16), Background: termenv.ANSI256Color(240)})
	segmentTime.SetContent(time.Now().Format(time.TimeOnly))
	segmentCursor.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(250), Background: termenv.ANSI256Color(237)})
	segmentCursor.SetContent("[1, 1]")

	style := powerline.StylePatched
	style.Color = prompt.Color{
		Foreground: termenv.ANSI256Color(234),
		Background: termenv.ANSI256Color(234),
	}

	p := powerline.Powerline{}
	p.AutoAdjustWidth(true) // use the full terminal-width
	p.Append(&segmentHost)
	p.Append(&segmentUser)
	p.Append(&segmentDBName)
	p.Append(&segmentCmdNum)
	p.AppendRight(&segmentCursor)
	p.AppendRight(&segmentDBType)
	p.AppendRight(&segmentTime)
	p.SetStyle(style)
	return p.Render
}
