package powerline

import (
	"testing"

	"github.com/jedib0t/go-prompter/prompt"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func TestSegment_Color(t *testing.T) {
	s := Segment{}

	c := s.Color()
	assert.Equal(t, termenv.BackgroundColor(), c.Background)
	assert.Equal(t, termenv.ForegroundColor(), c.Foreground)

	s.contentColor = &prompt.Color{
		Foreground: termenv.ANSI256Color(7),
		Background: termenv.ANSI256Color(16),
	}
	c = s.Color()
	assert.Equal(t, *s.contentColor, c)

	s.color = &prompt.Color{
		Foreground: termenv.ANSI256Color(17),
		Background: termenv.ANSI256Color(116),
	}
	c = s.Color()
	assert.Equal(t, *s.color, c)
}

func TestSegment_HasChanges(t *testing.T) {
	s := Segment{}
	assert.False(t, s.HasChanges())

	s.hasChanges = true
	assert.True(t, s.HasChanges())
}

func TestSegment_Render(t *testing.T) {
	s := Segment{}
	assert.Equal(t, "", s.Render())

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

	s.setPaddingLeft("[")
	s.setPaddingRight("]")
	assert.Equal(t, s.color.Sprint("[🌐 0.0.0.0]"), s.Render())
}

func TestSegment_ResetColor(t *testing.T) {
	s := Segment{}
	s.SetColor(prompt.Color{
		Foreground: termenv.ANSI256Color(0),
		Background: termenv.ANSI256Color(111),
	})
	assert.NotNil(t, s.color)

	s.ResetColor()
	assert.Nil(t, s.color)
}

func TestSegment_SetColor(t *testing.T) {
	s := Segment{}
	s.SetContent("foo")
	s.hasChanges = false
	assert.Nil(t, s.color)

	s.SetColor(prompt.Color{
		Foreground: termenv.ANSI256Color(0),
		Background: termenv.ANSI256Color(111),
	})
	assert.NotNil(t, s.color)
	assert.True(t, s.hasChanges)
	s.hasChanges = false

	s.SetColor(prompt.Color{
		Foreground: termenv.ANSI256Color(0),
		Background: termenv.ANSI256Color(111),
	})
	assert.NotNil(t, s.color)
	assert.False(t, s.hasChanges)
	s.hasChanges = false

	s.SetColor(prompt.Color{
		Foreground: termenv.ANSI256Color(100),
		Background: termenv.ANSI256Color(111),
	})
	assert.NotNil(t, s.color)
	assert.True(t, s.hasChanges)
}

func TestSegment_SetContent(t *testing.T) {
	s := Segment{}

	s.SetContent("foo", "tag1")
	assert.Equal(t, "foo", s.content)
	assert.NotNil(t, s.contentColor)
	assert.True(t, s.HasChanges())
	s.hasChanges = false

	s.SetContent("foo", "tag1")
	assert.Equal(t, "foo", s.content)
	assert.NotNil(t, s.contentColor)
	assert.False(t, s.HasChanges())
}

func TestSegment_SetIcon(t *testing.T) {
	s := Segment{}
	s.SetContent("foo")
	s.hasChanges = false
	assert.Empty(t, s.icon)

	s.SetIcon("foo")
	assert.Equal(t, "foo", s.icon)
	assert.True(t, s.hasChanges)
	s.hasChanges = false

	s.SetIcon("foo")
	assert.Equal(t, "foo", s.icon)
	assert.False(t, s.hasChanges)
}

func TestSegment_Width(t *testing.T) {
	s := Segment{}
	assert.Equal(t, 0, s.Width())

	s.SetContent("foo")
	assert.Equal(t, 6, s.Width())

	s.SetIcon("bar")
	assert.Equal(t, 10, s.Width())
}
