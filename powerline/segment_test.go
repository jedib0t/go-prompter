package powerline

import (
	"testing"

	"github.com/jedib0t/go-prompter/prompt"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func TestSegment_Render(t *testing.T) {
	s := Segment{}

	s.SetContent(testIP)
	assert.True(t, s.hasChanges)
	s.contentColor.Sprint(" " + testIP + " ")
	assert.Equal(t, s.contentColor.Sprint(" "+testIP+" "), s.Render())
	assert.False(t, s.hasChanges)

	s.SetIcon("🌐")
	assert.True(t, s.hasChanges)
	assert.Equal(t, s.contentColor.Sprint(" "+s.icon+" "+testIP+" "), s.Render())
	assert.False(t, s.hasChanges)

	s.SetColor(prompt.Color{
		Foreground: termenv.ANSI256Color(0),
		Background: termenv.ANSI256Color(111),
	})
	assert.True(t, s.hasChanges)
	assert.Equal(t, s.color.Sprint(" "+s.icon+" "+testIP+" "), s.Render())
	assert.False(t, s.hasChanges)
}
