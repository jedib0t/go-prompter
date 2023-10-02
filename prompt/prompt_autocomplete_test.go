package prompt

import (
	"testing"

	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func Test_printSuggestions(t *testing.T) {
	colorText := Color{Foreground: termenv.ANSI256Color(12), Background: termenv.ANSI256Color(13)}

	assert.Equal(t, "", printSuggestion("foo", colorText, 0))
	assert.Equal(t, colorText.Sprint(" ~ "), printSuggestion("foo", colorText, 1))
	assert.Equal(t, colorText.Sprint(" f~ "), printSuggestion("foo", colorText, 2))
	assert.Equal(t, colorText.Sprint(" foo "), printSuggestion("foo", colorText, 3))
	assert.Equal(t, colorText.Sprint(" foo  "), printSuggestion("foo", colorText, 4))
	assert.Equal(t, colorText.Sprint(" foo   "), printSuggestion("foo", colorText, 5))
}

func Test_printSuggestionsDropDown(t *testing.T) {
	suggestions := []Suggestion{
		{Value: "a", Hint: "A"},
		{Value: "bc", Hint: "B C"},
		{Value: "def", Hint: "D E F"},
		{Value: "ghij", Hint: "G H I J"},
		{Value: "klmno", Hint: "K L M N O"},
	}
	style := StyleAutoCompleteDefault
	style.NumItems = 3

	for _, idx := range []int{-1, 0} {
		output := printSuggestionsDropDown(suggestions, idx, style)
		expectedLines := []string{
			"\x1b[38;5;16;48;5;214m a        \x1b[0m\x1b[38;5;16;48;5;208m A         \x1b[0m\x1b[38;5;27;48;5;39m█\x1b[0m",
			"\x1b[38;5;16;48;5;45m bc       \x1b[0m\x1b[38;5;0;48;5;39m B C       \x1b[0m\x1b[38;5;27;48;5;39m░\x1b[0m",
			"\x1b[38;5;16;48;5;45m def      \x1b[0m\x1b[38;5;0;48;5;39m D E F     \x1b[0m\x1b[38;5;27;48;5;39m░\x1b[0m",
		}
		compareLines(t, expectedLines, output, idx)
	}

	output := printSuggestionsDropDown(suggestions, 2, style)
	expectedLines := []string{
		"\x1b[38;5;16;48;5;45m bc       \x1b[0m\x1b[38;5;0;48;5;39m B C       \x1b[0m\x1b[38;5;27;48;5;39m░\x1b[0m",
		"\x1b[38;5;16;48;5;214m def      \x1b[0m\x1b[38;5;16;48;5;208m D E F     \x1b[0m\x1b[38;5;27;48;5;39m█\x1b[0m",
		"\x1b[38;5;16;48;5;45m ghij     \x1b[0m\x1b[38;5;0;48;5;39m G H I J   \x1b[0m\x1b[38;5;27;48;5;39m░\x1b[0m",
	}
	compareLines(t, expectedLines, output, 2)

	for _, idx := range []int{4, 8} {
		output = printSuggestionsDropDown(suggestions, idx, style)
		expectedLines = []string{
			"\x1b[38;5;16;48;5;45m def      \x1b[0m\x1b[38;5;0;48;5;39m D E F     \x1b[0m\x1b[38;5;27;48;5;39m░\x1b[0m",
			"\x1b[38;5;16;48;5;45m ghij     \x1b[0m\x1b[38;5;0;48;5;39m G H I J   \x1b[0m\x1b[38;5;27;48;5;39m░\x1b[0m",
			"\x1b[38;5;16;48;5;214m klmno    \x1b[0m\x1b[38;5;16;48;5;208m K L M N O \x1b[0m\x1b[38;5;27;48;5;39m█\x1b[0m",
		}
		compareLines(t, expectedLines, output, idx)
	}
}
