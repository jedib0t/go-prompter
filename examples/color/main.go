package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jedib0t/go-prompter/prompt"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

var (
	colorsTitle = prompt.Color{
		Foreground: termenv.BackgroundColor(),
		Background: termenv.ForegroundColor(),
	}
)

//gocyclo:ignore
func main() {
	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	out := strings.Builder{}
	out.WriteString("\n")
	out.WriteString(colorsTitle.Sprint(" 256-color Mode "))
	out.WriteString("\n")
	for bg := 0; bg < 256; bg++ {
		if bg == 0 {
			out.WriteString("\n      Standard Colors  High-Intensity Colors\n")
		} else if bg == 8 {
			out.WriteString("  ")
		} else if bg == 16 {
			out.WriteString("\n\n216 Colors\n")
		} else if bg > 16 && bg < 232 && (bg-16)%18 == 0 { // 216-colors
			if termWidth < 180 {
				out.WriteString("\n") // 18 colors in one line
			} else if (bg-16)%36 == 0 {
				out.WriteString("\n") // 36 colors in one line
			}
		} else if bg == 232 {
			out.WriteString("\n\nGrayscale Colors\n")
		} else if bg == 244 {
			if termWidth < 120 { // 12 grayscale colors in one line
				out.WriteString("\n")
			}
		}

		// determine foreground color based on how bright the background is
		fg := 16 // black by default
		if bg >= 16 && (bg-16)%36 < 18 && bg < 244 {
			fg = 15 // white for darker background
		}
		c := prompt.Color{
			Foreground: termenv.ANSI256Color(fg),
			Background: termenv.ANSI256Color(bg),
		}
		out.WriteString(c.Sprintf(" %03d ", bg))
	}
	out.WriteString("\n")

	for _, line := range strings.Split(out.String(), "\n") {
		fmt.Println(text.AlignCenter.Apply(line, termWidth))
	}
}
