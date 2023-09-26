package prompt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func BenchmarkWidthEnforcerDefault(b *testing.B) {
	color := Color{
		Background: termenv.ANSI256Color(100),
		Foreground: termenv.ANSI256Color(200),
	}
	in := color.Sprint("Ghosts of the Deep")
	out := WidthEnforcerDefault(in, 5)
	expected := "" +
		"\x1b[38;5;200;48;5;100mGhost\x1b[0m\n" +
		"\x1b[38;5;200;48;5;100ms of \x1b[0m\n" +
		"\x1b[38;5;200;48;5;100mthe D\x1b[0m\n" +
		"\x1b[38;5;200;48;5;100meep\x1b[0m"
	assert.Equal(b, expected, out)

	for idx := 0; idx < b.N; idx++ {
		WidthEnforcerDefault(in, 5)
	}
}

func TestWidthEnforcerDefault(t *testing.T) {
	in := "Ghosts of the Deep"
	out := WidthEnforcerDefault(in, 5)
	expected := "Ghost\ns of \nthe D\neep"
	assert.Equal(t, expected, out)

	renderTestExpectedCode := func() {
		outLines := strings.Split(out, "\n")
		for idx, line := range outLines {
			if idx == 0 {
				fmt.Printf("    expected = \"\" +\n")
			} else {
				fmt.Print(" +\n")
			}
			if idx < len(outLines)-1 {
				line += "\n"
			}
			fmt.Printf("        %#v", line)
		}
		fmt.Printf("\n")
	}

	color := Color{
		Background: termenv.ANSI256Color(100),
		Foreground: termenv.ANSI256Color(200),
	}
	in = color.Sprint("Ghosts of the Deep")
	out = WidthEnforcerDefault(in, 5)
	expected = "" +
		"\x1b[38;5;200;48;5;100mGhost\x1b[0m\n" +
		"\x1b[38;5;200;48;5;100ms of \x1b[0m\n" +
		"\x1b[38;5;200;48;5;100mthe D\x1b[0m\n" +
		"\x1b[38;5;200;48;5;100meep\x1b[0m"
	assert.Equal(t, expected, out)
	if expected != out {
		renderTestExpectedCode()
	}

	out = WidthEnforcerDefault(in, 8)
	expected = "" +
		"\x1b[38;5;200;48;5;100mGhosts o\x1b[0m\n" +
		"\x1b[38;5;200;48;5;100mf the De\x1b[0m\n" +
		"\x1b[38;5;200;48;5;100mep\x1b[0m"
	assert.Equal(t, expected, out)
	if expected != out {
		renderTestExpectedCode()
	}
}
