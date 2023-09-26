package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/user"
	"time"

	"github.com/jedib0t/go-prompter/powerline"
	"github.com/jedib0t/go-prompter/prompt"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

var (
	flagIP = flag.String("ip", "", "Use this IP Address instead of looking up localhost's IP")
)

func main() {
	flag.Parse()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	userObj, err := user.Current()
	if err != nil {
		userObj = &user.User{Username: "username"}
	}
	username := userObj.Username

	hostIP := "localhost"
	if *flagIP != "" {
		hostIP = *flagIP
	} else if conn, err := net.Dial("udp", "8.8.8.8:80"); err == nil {
		defer func() {
			_ = conn.Close()
		}()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		hostIP = localAddr.IP.String()
	}

	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	if termWidth < 0 {
		termWidth = 80
	}

	header := prompt.LineRuler()

	// A Powerline is made up of multiple segments stitched together. First,
	// prepare all the segments that are going to be used.
	segWidth := powerline.Segment{}
	segHostname := powerline.Segment{}
	segHostname.SetContent(hostname)
	segHostname.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(23)})
	segHostname.SetIcon("💻")
	segUsername := powerline.Segment{}
	segUsername.SetContent(username)
	segUsername.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(17)})
	segUsername.SetIcon("👤")
	// For branch, we want to use a common color across all branches but don't
	// want to call SetColor manually. In this case, the SetContent call uses
	// the tag "git-branch" for determining color instead of the content "main".
	segBranch := powerline.Segment{}
	segBranch.SetContent("main", "git-branch")
	segBranch.SetIcon("")
	segCmdNum := powerline.Segment{}
	segCmdNum.SetContent("1")
	segCmdNum.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(16), Background: termenv.ANSI256Color(147)})
	segHostIP := powerline.Segment{}
	segHostIP.SetContent(hostIP)
	segHostIP.SetIcon("🌐")
	segTime := powerline.Segment{}
	segTime.SetColor(prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(239)})

	style := powerline.StylePatched
	style.Color = prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.BackgroundColor()}
	fmt.Println("Simple Powerline with segments on the left:")
	fmt.Println(header(termWidth))
	p := powerline.Powerline{}
	p.Append(&segWidth)
	p.Append(&segHostname)
	p.Append(&segUsername)
	p.Append(&segBranch)
	p.Append(&segCmdNum)
	p.SetStyle(style)
	for _, width := range []int{0, 20, 40, 80, 120, termWidth} {
		segWidth.SetContent(fmt.Sprintf("Width: %03d", width))
		fmt.Printf("%s\n", p.Render(width))
	}
	fmt.Println()

	style = powerline.StylePatched
	style.Color = prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(234)}
	fmt.Println("Powerline with segments on both sides:")
	fmt.Println(header(termWidth))
	p = powerline.Powerline{}
	p.Append(&segWidth)
	p.Append(&segHostname)
	p.Append(&segUsername)
	p.Append(&segBranch)
	p.Append(&segCmdNum)
	p.AppendRight(&segHostIP)
	p.AppendRight(&segTime)
	p.SetStyle(style)
	segTime.SetContent(time.Now().Format("15:04:05.000"))
	for _, width := range []int{0, 20, 40, 80, 120, termWidth} {
		segWidth.SetContent(fmt.Sprintf("Width: %03d", width))
		fmt.Printf("%s\n", p.Render(width))
	}
	fmt.Println()

	style = powerline.StyleNonPatched
	style.Color = prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(234)}
	fmt.Println("Powerline with segments on both sides using non-patched separators:")
	fmt.Println(header(termWidth))
	p = powerline.Powerline{}
	p.Append(&segWidth)
	p.Append(&segHostname)
	p.Append(&segUsername)
	p.Append(&segBranch)
	p.Append(&segCmdNum)
	p.AppendRight(&segHostIP)
	p.AppendRight(&segTime)
	p.SetStyle(style)
	segTime.SetContent(time.Now().Format("15:04:05.000"))
	for _, width := range []int{0, 20, 40, 80, 120, termWidth} {
		segWidth.SetContent(fmt.Sprintf("Width: %03d", width))
		fmt.Printf("%s\n", p.Render(width))
	}
	fmt.Println()

	style = powerline.StylePatched
	style.Color = prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(234)}
	fmt.Println("Powerline with segments switched between sides:")
	fmt.Println(header(termWidth))
	p = powerline.Powerline{}
	p.Append(&segWidth)
	p.Append(&segTime)
	p.Append(&segHostIP)
	p.AppendRight(&segCmdNum)
	p.AppendRight(&segBranch)
	p.AppendRight(&segUsername)
	p.AppendRight(&segHostname)
	p.SetStyle(style)
	segTime.SetContent(time.Now().Format("15:04:05.000"))
	for _, width := range []int{0, 20, 40, 80, 120, termWidth} {
		segWidth.SetContent(fmt.Sprintf("Width: %03d", width))
		fmt.Printf("%s\n", p.Render(width))
	}
	fmt.Println()

	style = powerline.StylePatched
	style.Color = prompt.Color{Foreground: termenv.ANSI256Color(7), Background: termenv.ANSI256Color(234)}
	fmt.Println("Powerline with segments on both sides with auto-adjusting width:")
	fmt.Println(header(termWidth))
	p = powerline.Powerline{}
	p.AutoAdjustWidth(true)
	p.Append(&segWidth)
	p.Append(&segHostname)
	p.Append(&segUsername)
	p.Append(&segBranch)
	p.Append(&segCmdNum)
	p.AppendRight(&segHostIP)
	p.AppendRight(&segTime)
	p.SetStyle(style)
	segTime.SetContent(time.Now().Format("15:04:05.000"))
	for _, width := range []int{0, 20, 40, 80, 120, termWidth} {
		segWidth.SetContent(fmt.Sprintf("Width: %03d", width))
		fmt.Printf("%s\n", p.Render(width))
	}
	fmt.Println()
}
