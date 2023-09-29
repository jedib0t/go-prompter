package prompt

import (
	"testing"

	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func TestColor_Invert(t *testing.T) {
	c := Color{
		Foreground: termenv.ANSI256Color(194),
		Background: termenv.ANSI256Color(56),
	}
	c2 := c.Invert()

	assert.Equal(t, termenv.ANSI256Color(194), c.Foreground)
	assert.Equal(t, termenv.ANSI256Color(56), c.Background)
	assert.Equal(t, termenv.ANSI256Color(56), c2.Foreground)
	assert.Equal(t, termenv.ANSI256Color(194), c2.Background)
}

func TestColor_Sprint(t *testing.T) {
	c := Color{}
	assert.Equal(t, "foo", c.Sprint("foo"))

	c = Color{
		Foreground: termenv.ANSI256Color(194),
		Background: termenv.ANSI256Color(56),
	}
	assert.Equal(t, "\x1b[38;5;194;48;5;56mfoo\x1b[0m", c.Sprint("foo"))
}

func TestColor_Sprintf(t *testing.T) {
	c := Color{
		Foreground: termenv.ANSI256Color(194),
		Background: termenv.ANSI256Color(56),
	}
	assert.Equal(t, "\x1b[38;5;194;48;5;56mfoo   bar\x1b[0m", c.Sprintf("foo %5s", "bar"))
}
