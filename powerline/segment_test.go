package powerline

import (
	"testing"

	"github.com/jedib0t/go-prompter/prompt"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func TestSegment_Render(t *testing.T) {
	s := Segment{}

	s.SetContent("10.0.0.1")
	assert.True(t, s.hasChanges)
	assert.Equal(t, "\x1b[38;5;16;48;5;188m 10.0.0.1 \x1b[0m", s.Render())
	assert.False(t, s.hasChanges)

	s.SetIcon("🌐")
	assert.True(t, s.hasChanges)
	assert.Equal(t, "\x1b[38;5;16;48;5;188m 🌐 10.0.0.1 \x1b[0m", s.Render())
	assert.False(t, s.hasChanges)

	s.SetColor(prompt.Color{
		Foreground: termenv.ANSI256Color(0),
		Background: termenv.ANSI256Color(111),
	})
	assert.True(t, s.hasChanges)
	assert.Equal(t, "\x1b[38;5;0;48;5;111m 🌐 10.0.0.1 \x1b[0m", s.Render())
	assert.False(t, s.hasChanges)
}
