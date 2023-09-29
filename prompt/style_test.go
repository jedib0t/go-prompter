package prompt

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStyleValidate(t *testing.T) {
	s := StyleDefault
	err := s.Validate()
	assert.Nil(t, err)

	s = StyleDefault
	s.Dimensions.HeightMin = 5
	s.Dimensions.HeightMax = 4
	err = s.Validate()
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrInvalidDimensions))
	assert.Contains(t, err.Error(), "height-min [5] cannot be greater than height-max [4]")

	s = StyleDefault
	s.Dimensions.WidthMin = 50
	s.Dimensions.WidthMax = 40
	err = s.Validate()
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrInvalidDimensions))
	assert.Contains(t, err.Error(), "width-min [50] cannot be greater than width-max [40]")
}

func TestScrollbarGenerate(t *testing.T) {
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

	expectedIndicatorEmpty := s.Color.Sprint(fmt.Sprintf("%c", s.IndicatorEmpty))
	expectedIndicator := s.Color.Sprint(fmt.Sprintf("%c", s.Indicator))

	expectedLines = []string{
		expectedIndicator,
		expectedIndicatorEmpty,
		expectedIndicatorEmpty,
		expectedIndicatorEmpty,
		expectedIndicatorEmpty,
	}
	actualLines, isVisible = s.Generate(20, 0, 5)
	compareModelLines(t, expectedLines, actualLines, "Scrollbar @ 0%")
	assert.True(t, isVisible)

	expectedLines = []string{
		expectedIndicatorEmpty,
		expectedIndicator,
		expectedIndicatorEmpty,
		expectedIndicatorEmpty,
		expectedIndicatorEmpty,
	}
	actualLines, isVisible = s.Generate(20, 5, 5)
	compareModelLines(t, expectedLines, actualLines, "Scrollbar @ 25%")
	assert.True(t, isVisible)

	expectedLines = []string{
		expectedIndicatorEmpty,
		expectedIndicatorEmpty,
		expectedIndicator,
		expectedIndicatorEmpty,
		expectedIndicatorEmpty,
	}
	actualLines, isVisible = s.Generate(20, 10, 5)
	compareModelLines(t, expectedLines, actualLines, "Scrollbar @ 50%")
	assert.True(t, isVisible)

	expectedLines = []string{
		expectedIndicatorEmpty,
		expectedIndicatorEmpty,
		expectedIndicatorEmpty,
		expectedIndicator,
		expectedIndicatorEmpty,
	}
	actualLines, isVisible = s.Generate(20, 15, 5)
	compareModelLines(t, expectedLines, actualLines, "Scrollbar @ 75%")
	assert.True(t, isVisible)

	expectedLines = []string{
		expectedIndicatorEmpty,
		expectedIndicatorEmpty,
		expectedIndicatorEmpty,
		expectedIndicatorEmpty,
		expectedIndicator,
	}
	actualLines, isVisible = s.Generate(20, 20, 5)
	compareModelLines(t, expectedLines, actualLines, "Scrollbar @ 100%")
	assert.True(t, isVisible)
}
