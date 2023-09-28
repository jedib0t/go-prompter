package powerline

import (
	"testing"
	"time"

	"github.com/jedib0t/go-prompter/prompt"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

var (
	testIP = "0.0.0.0"
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
	segHostIP.SetContent(testIP)
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
	segHostIP.SetContent(testIP)
	segHostIP.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(16), Background: termenv.ANSI256Color(201)})
	segTime := &Segment{}
	segTime.SetContent("12:13:14")
	segTime.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(239)})

	style := StyleNonPatched
	style.Color = prompt.Color{Foreground: termenv.ANSI256Color(235), Background: termenv.ANSI256Color(235)}

	t.Run("with segments on one side", func(t *testing.T) {
		p := Powerline{}
		p.Append(segUser)
		p.Append(segHost)
		p.Append(segCmdNum)

		assert.Equal(t,
			segUser.Render()+segHost.Render()+segCmdNum.Render(),
			p.Render(0),
		)
	})

	expectedOut120 := segUser.Render() +
		"\x1b[38;5;205;48;5;17m" + style.SeparatorLeft + "\x1b[0m" +
		segHost.Render() +
		"\x1b[38;5;147;48;5;205m" + style.SeparatorLeft + "\x1b[0m" +
		segCmdNum.Render() +
		"\x1b[38;5;235;48;5;147m" + style.SeparatorLeft + "\x1b[0m" +
		style.Color.Sprint("                                                                      ") +
		"\x1b[38;5;235;48;5;201m" + style.SeparatorRight + "\x1b[0m" +
		segHostIP.Render() +
		"\x1b[38;5;201;48;5;239m" + style.SeparatorRight + "\x1b[0m" +
		segTime.Render()

	expectedOut50 := segUser.Render() +
		"\x1b[38;5;205;48;5;17m" + style.SeparatorLeft + "\x1b[0m" +
		segHost.Render() +
		"\x1b[38;5;147;48;5;205m" + style.SeparatorLeft + "\x1b[0m" +
		segCmdNum.Render() +
		"\x1b[38;5;235;48;5;147m" + style.SeparatorLeft + "\x1b[0m" +
		style.Color.Sprint("") +
		"\x1b[38;5;235;48;5;201m" + style.SeparatorRight + "\x1b[0m" +
		segHostIP.Render() +
		"\x1b[38;5;201;48;5;239m" + style.SeparatorRight + "\x1b[0m" +
		segTime.Render()

	t.Run("without auto-adjusting width", func(t *testing.T) {
		p := Powerline{}
		p.Append(segUser)
		p.Append(segHost)
		p.Append(segCmdNum)
		p.AppendRight(segHostIP)
		p.AppendRight(segTime)
		p.SetStyle(style)
		p.AutoAdjustWidth(false)

		assert.Equal(t, expectedOut120, p.Render(120))
		assert.Equal(t, expectedOut50, p.Render(50))
		assert.Equal(t,
			segUser.Render()+
				"\x1b[38;5;205;48;5;17m"+style.SeparatorLeft+"\x1b[0m"+
				segHost.Render()+
				"\x1b[38;5;147;48;5;205m"+style.SeparatorLeft+"\x1b[0m"+
				segCmdNum.Render()+
				"\x1b[38;5;235;48;5;147m"+style.SeparatorLeft+"\x1b[0m"+
				style.Color.Sprint("")+
				"\x1b[38;5;235;48;5;201m"+style.SeparatorRight+"\x1b[0m"+
				segHostIP.Render()+
				"\x1b[38;5;201;48;5;239m"+style.SeparatorRight+"\x1b[0m"+
				segTime.Render(),
			p.Render(25),
		)
	})

	t.Run("with auto-adjusting width", func(t *testing.T) {
		p := Powerline{}
		p.Append(segUser)
		p.Append(segHost)
		p.Append(segCmdNum)
		p.AppendRight(segHostIP)
		p.AppendRight(segTime)
		p.SetStyle(style)
		p.AutoAdjustWidth(true)

		assert.Equal(t, expectedOut120, p.Render(120))
		assert.Equal(t, expectedOut50, p.Render(50))
		assert.Equal(t, expectedOut50, p.Render(50)) // to check cache
		assert.Equal(t,
			segUser.Render()+
				"\x1b[38;5;235;48;5;17m"+style.SeparatorLeft+"\x1b[0m"+
				style.Color.Sprint("")+
				"\x1b[38;5;235;48;5;239m"+style.SeparatorRight+"\x1b[0m"+
				segTime.Render(),
			p.Render(25),
		)
	})
}
