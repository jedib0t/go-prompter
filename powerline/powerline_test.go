package powerline

import (
	"testing"
	"time"

	"github.com/jedib0t/go-prompter/prompt"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func BenchmarkPowerline_Render(b *testing.B) {
	segHostname := &Segment{}
	segHostname.SetContent("hostname")
	segHostname.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(23)})
	segHostname.SetIcon("💻")
	segUsername := &Segment{}
	segUsername.SetContent("username")
	segUsername.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(17)})
	segUsername.SetIcon("👤")
	segBranch := &Segment{}
	segBranch.SetContent("main", "git-branch")
	segBranch.SetIcon("")
	segCmdNum := &Segment{}
	segCmdNum.SetContent("1")
	segCmdNum.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(16), Background: termenv.ANSI256Color(147)})
	segTime := &Segment{}
	segTime.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(239)})
	segHostIP := &Segment{}
	segHostIP.SetContent("10.0.0.10")
	segHostIP.SetIcon("🌐")

	style := StyleDefault
	style.Color = prompt.Color{Foreground: termenv.ANSI256Color(235), Background: termenv.ANSI256Color(235)}

	p := Powerline{}
	p.Append(segHostname)
	p.Append(segUsername)
	p.Append(segBranch)
	p.Append(segCmdNum)
	p.AppendRight(segHostIP)
	p.AppendRight(segTime)
	p.SetStyle(style)
	for i := 0; i < b.N; i++ {
		segTime.SetContent(time.Now().Format("15:04:05.000000"))
		p.Render(120)
	}
}

func TestPowerline_Render(t *testing.T) {
	segUser := &Segment{}
	segUser.SetContent("username")
	segUser.SetIcon("👤")
	segUser.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(17)})
	segHost := &Segment{}
	segHost.SetContent("hostname")
	segCmdNum := &Segment{}
	segCmdNum.SetContent("1")
	segCmdNum.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(16), Background: termenv.ANSI256Color(147)})
	segHostIP := &Segment{}
	segHostIP.SetContent("10.0.0.1")
	segHostIP.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(16), Background: termenv.ANSI256Color(201)})
	segTime := &Segment{}
	segTime.SetContent("12:13:14")
	segTime.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(239)})

	style := StyleNonPatched
	style.Color = prompt.Color{Foreground: termenv.ANSI256Color(235), Background: termenv.ANSI256Color(235)}

	p := Powerline{}
	p.Append(segUser)
	p.Append(segHost)
	p.Append(segCmdNum)
	p.AppendRight(segHostIP)
	p.AppendRight(segTime)
	p.SetStyle(style)
	assert.Equal(t,
		segUser.Color().Sprint(" 👤 username ")+
			"\x1b[38;5;205;48;5;17m"+style.SeparatorLeft+"\x1b[0m"+
			segHost.Color().Sprint(" hostname ")+
			"\x1b[38;5;147;48;5;205m"+style.SeparatorLeft+"\x1b[0m"+
			segCmdNum.Color().Sprint(" 1 ")+
			"\x1b[38;5;235;48;5;147m"+style.SeparatorLeft+"\x1b[0m"+
			"\x1b[38;5;235;48;5;235m                                                                     \x1b[0m"+
			"\x1b[38;5;235;48;5;201m"+style.SeparatorRight+"\x1b[0m"+
			segHostIP.Color().Sprint(" 10.0.0.1 ")+
			"\x1b[38;5;201;48;5;239m"+style.SeparatorRight+"\x1b[0m"+
			segTime.Color().Sprint(" 12:13:14 "),
		p.Render(120))
	assert.Equal(t,
		segUser.Color().Sprint(" 👤 username ")+
			"\x1b[38;5;205;48;5;17m"+style.SeparatorLeft+"\x1b[0m"+
			segHost.Color().Sprint(" hostname ")+
			"\x1b[38;5;147;48;5;205m"+style.SeparatorLeft+"\x1b[0m"+
			segCmdNum.Color().Sprint(" 1 ")+
			"\x1b[38;5;235;48;5;147m"+style.SeparatorLeft+"\x1b[0m"+
			"\x1b[38;5;235;48;5;235m\x1b[0m"+
			"\x1b[38;5;235;48;5;201m"+style.SeparatorRight+"\x1b[0m"+
			segHostIP.Color().Sprint(" 10.0.0.1 ")+
			"\x1b[38;5;201;48;5;239m"+style.SeparatorRight+"\x1b[0m"+
			segTime.Color().Sprint(" 12:13:14 "),
		p.Render(50))
}

func TestPowerline_RenderWithAutoAdjust(t *testing.T) {
	segUser := &Segment{}
	segUser.SetContent("username")
	segUser.SetIcon("👤")
	segUser.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(17)})
	segHost := &Segment{}
	segHost.SetContent("hostname")
	segCmdNum := &Segment{}
	segCmdNum.SetContent("1")
	segCmdNum.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(16), Background: termenv.ANSI256Color(147)})
	segHostIP := &Segment{}
	segHostIP.SetContent("10.0.0.1")
	segHostIP.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(16), Background: termenv.ANSI256Color(201)})
	segTime := &Segment{}
	segTime.SetContent("12:13:14")
	segTime.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(239)})

	style := StyleNonPatched
	style.Color = prompt.Color{Foreground: termenv.ANSI256Color(235), Background: termenv.ANSI256Color(235)}

	p := Powerline{}
	p.AutoAdjustWidth(true)
	p.Append(segUser)
	p.Append(segHost)
	p.Append(segCmdNum)
	p.AppendRight(segHostIP)
	p.AppendRight(segTime)
	p.SetStyle(style)
	assert.Equal(t,
		segUser.Color().Sprint(" 👤 username ")+
			"\x1b[38;5;205;48;5;17m"+style.SeparatorLeft+"\x1b[0m"+
			segHost.Color().Sprint(" hostname ")+
			"\x1b[38;5;147;48;5;205m"+style.SeparatorLeft+"\x1b[0m"+
			segCmdNum.Color().Sprint(" 1 ")+
			"\x1b[38;5;235;48;5;147m"+style.SeparatorLeft+"\x1b[0m"+
			"\x1b[38;5;235;48;5;235m                                                                     \x1b[0m"+
			"\x1b[38;5;235;48;5;201m"+style.SeparatorRight+"\x1b[0m"+
			segHostIP.Color().Sprint(" 10.0.0.1 ")+
			"\x1b[38;5;201;48;5;239m"+style.SeparatorRight+"\x1b[0m"+
			segTime.Color().Sprint(" 12:13:14 "),
		p.Render(120))
	assert.Equal(t,
		segUser.Color().Sprint(" 👤 username ")+
			"\x1b[38;5;205;48;5;17m"+style.SeparatorLeft+"\x1b[0m"+
			segHost.Color().Sprint(" hostname ")+
			"\x1b[38;5;235;48;5;205m"+style.SeparatorLeft+"\x1b[0m"+
			"\x1b[38;5;235;48;5;235m   \x1b[0m"+
			"\x1b[38;5;235;48;5;201m"+style.SeparatorRight+"\x1b[0m"+
			segHostIP.Color().Sprint(" 10.0.0.1 ")+
			"\x1b[38;5;201;48;5;239m"+style.SeparatorRight+"\x1b[0m"+
			segTime.Color().Sprint(" 12:13:14 "),
		p.Render(50))
	assert.Equal(t,
		segUser.Color().Sprint(" 👤 username ")+
			"\x1b[38;5;235;48;5;17m"+style.SeparatorLeft+"\x1b[0m"+
			"\x1b[38;5;235;48;5;235m\x1b[0m"+
			"\x1b[38;5;235;48;5;239m"+style.SeparatorRight+"\x1b[0m"+
			segTime.Color().Sprint(" 12:13:14 "),
		p.Render(25))
}
