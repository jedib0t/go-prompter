package powerline

import (
	"testing"

	"github.com/jedib0t/go-prompter/prompt"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func TestSegment_Render(t *testing.T) {
	s := Segment{}

	s.SetContent("0.0.0.0")
	assert.True(t, s.hasChanges)
	assert.Equal(t, "\x1b[38;5;7;48;5;91m 0.0.0.0 \x1b[0m", s.Render())
	assert.False(t, s.hasChanges)

	s.SetIcon("🌐")
	assert.True(t, s.hasChanges)
	assert.Equal(t, "\x1b[38;5;7;48;5;91m 🌐 0.0.0.0 \x1b[0m", s.Render())
	assert.False(t, s.hasChanges)

	s.SetColor(prompt.Color{
		Foreground: termenv.ANSI256Color(0),
		Background: termenv.ANSI256Color(111),
	})
	assert.True(t, s.hasChanges)
	assert.Equal(t, "\x1b[38;5;0;48;5;111m 🌐 0.0.0.0 \x1b[0m", s.Render())
	assert.False(t, s.hasChanges)
}
