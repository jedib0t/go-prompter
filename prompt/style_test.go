package prompt

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStyle_validate(t *testing.T) {
	s := StyleDefault
	err := s.validate()
	assert.Nil(t, err)

	s = StyleDefault
	s.Dimensions.HeightMin = 5
	s.Dimensions.HeightMax = 4
	err = s.validate()
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrInvalidDimensions))
	assert.Contains(t, err.Error(), "height-min [5] cannot be greater than height-max [4]")

	s = StyleDefault
	s.Dimensions.WidthMin = 50
	s.Dimensions.WidthMax = 40
	err = s.validate()
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrInvalidDimensions))
	assert.Contains(t, err.Error(), "width-min [50] cannot be greater than width-max [40]")
}

func TestScrollbar_Generate(t *testing.T) {
	s := StyleScrollbarDefault

	expectedLines := []string{
		"",
		"",
		"",
		"",
		"",
	}
	actualLines, isVisible := s.Generate(5, 0, 5)
	compareModelLines(t, expectedLines, actualLines, "Scrollbar None")
	assert.False(t, isVisible)

	expectedLines = []string{
		"\x1b[38;5;237;48;5;233m█\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
	}
	actualLines, isVisible = s.Generate(20, 0, 5)
	compareModelLines(t, expectedLines, actualLines, "Scrollbar @ 0%")
	assert.True(t, isVisible)

	expectedLines = []string{
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m█\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
	}
	actualLines, isVisible = s.Generate(20, 5, 5)
	compareModelLines(t, expectedLines, actualLines, "Scrollbar @ 25%")
	assert.True(t, isVisible)

	expectedLines = []string{
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m█\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
	}
	actualLines, isVisible = s.Generate(20, 10, 5)
	compareModelLines(t, expectedLines, actualLines, "Scrollbar @ 50%")
	assert.True(t, isVisible)

	expectedLines = []string{
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m█\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
	}
	actualLines, isVisible = s.Generate(20, 15, 5)
	compareModelLines(t, expectedLines, actualLines, "Scrollbar @ 75%")
	assert.True(t, isVisible)

	expectedLines = []string{
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m░\x1b[0m",
		"\x1b[38;5;237;48;5;233m█\x1b[0m",
	}
	actualLines, isVisible = s.Generate(20, 20, 5)
	compareModelLines(t, expectedLines, actualLines, "Scrollbar @ 100%")
	assert.True(t, isVisible)
}
