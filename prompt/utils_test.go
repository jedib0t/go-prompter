package prompt

import (
	"fmt"
	"testing"

	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func Test_calculateViewportRange(t *testing.T) {
	start, stop := calculateViewportRange(5, -1, -1)
	assert.Equal(t, 0, start)
	assert.Equal(t, 4, stop)

	start, stop = calculateViewportRange(5, -1, 5)
	assert.Equal(t, 0, start)
	assert.Equal(t, 4, stop)

	start, stop = calculateViewportRange(5, 6, 5)
	assert.Equal(t, 0, start)
	assert.Equal(t, 4, stop)

	start, stop = calculateViewportRange(3, 6, 5)
	assert.Equal(t, 0, start)
	assert.Equal(t, 2, stop)

	start, stop = calculateViewportRange(10, 11, 5)
	assert.Equal(t, 5, start)
	assert.Equal(t, 9, stop)

	start, stop = calculateViewportRange(5, 0, 5)
	assert.Equal(t, 0, start)
	assert.Equal(t, 4, stop)

	start, stop = calculateViewportRange(5, 1, 5)
	assert.Equal(t, 0, start)
	assert.Equal(t, 4, stop)

	start, stop = calculateViewportRange(5, 2, 5)
	assert.Equal(t, 0, start)
	assert.Equal(t, 4, stop)

	start, stop = calculateViewportRange(5, 4, 5)
	assert.Equal(t, 0, start)
	assert.Equal(t, 4, stop)

	start, stop = calculateViewportRange(8, 4, 5)
	assert.Equal(t, 2, start)
	assert.Equal(t, 6, stop)

	start, stop = calculateViewportRange(8, 8, 5)
	assert.Equal(t, 3, start)
	assert.Equal(t, 7, stop)

	start, stop = calculateViewportRange(8, 9, 5)
	assert.Equal(t, 3, start)
	assert.Equal(t, 7, stop)
}

func Test_clampValue(t *testing.T) {
	assert.Equal(t, 5, clampValue(3, 5, 10))
	assert.Equal(t, 5, clampValue(5, 5, 10))
	assert.Equal(t, 7, clampValue(7, 5, 10))
	assert.Equal(t, 10, clampValue(10, 5, 10))
	assert.Equal(t, 10, clampValue(12, 5, 10))

	assert.Equal(t, 3, clampValue(3, 0, 0))
	assert.Equal(t, 5, clampValue(3, 5, 0))
	assert.Equal(t, 5, clampValue(6, 0, 5))
}

func Test_clampValueAllowZero(t *testing.T) {
	assert.Equal(t, 0, clampValueAllowZero(0, 5, 10))

	assert.Equal(t, 5, clampValueAllowZero(3, 5, 10))
	assert.Equal(t, 5, clampValueAllowZero(5, 5, 10))
	assert.Equal(t, 7, clampValueAllowZero(7, 5, 10))
	assert.Equal(t, 10, clampValueAllowZero(10, 5, 10))
	assert.Equal(t, 10, clampValueAllowZero(12, 5, 10))

	assert.Equal(t, 3, clampValueAllowZero(3, 0, 0))
	assert.Equal(t, 5, clampValueAllowZero(3, 5, 0))
	assert.Equal(t, 5, clampValueAllowZero(6, 0, 5))
}

func Benchmark_insertCursor(b *testing.B) {
	colorContent1 := Color{Foreground: termenv.ANSI256Color(81), Background: termenv.ANSI256Color(0)}
	colorCursor := StyleCursorDefault.Color

	input := colorContent1.Sprint("select")
	expectedOutput := colorCursor.Sprint("s") +
		colorContent1.Sprint("elect")
	output := insertCursor(input, 0, colorCursor)
	assert.Equal(b, expectedOutput, output)

	for idx := 0; idx < b.N; idx++ {
		insertCursor(input, 0, colorCursor)
	}
}

func Test_insertCursor(t *testing.T) {
	colorContent1 := Color{Foreground: termenv.ANSI256Color(81), Background: termenv.ANSI256Color(0)}
	colorContent2 := Color{Foreground: termenv.ANSI256Color(82), Background: termenv.ANSI256Color(0)}
	colorCursor := StyleCursorDefault.Color

	input := "select"
	expectedOutput := colorCursor.Sprint("s") + "elect"
	output := insertCursor(input, 0, colorCursor)
	assert.Equal(t, expectedOutput, output)

	input = colorContent1.Sprint("select")
	expectedOutput = colorCursor.Sprint("s") +
		colorContent1.Sprint("elect")
	output = insertCursor(input, 0, colorCursor)
	assert.Equal(t, expectedOutput, output)

	input = colorContent1.Sprint("select")
	expectedOutput = colorContent1.Sprint("select") +
		colorCursor.Sprint(" ")
	output = insertCursor(input, 10, colorCursor)
	assert.Equal(t, expectedOutput, output)

	input = colorContent1.Sprint("select") +
		colorContent2.Sprint(" ") +
		colorContent1.Sprint("foo")
	expectedOutput = colorContent1.Sprint("select") +
		colorContent2.Sprint(" ") +
		colorContent1.Sprint("f") +
		colorCursor.Sprint("o") +
		colorContent1.Sprint("o")
	output = insertCursor(input, 8, colorCursor)
	assert.Equal(t, expectedOutput, output)
}

func Test_overwriteContent(t *testing.T) {
	colorContent1 := Color{Foreground: termenv.ANSI256Color(0), Background: termenv.ANSI256Color(12)}
	colorContent2 := Color{Foreground: termenv.ANSI256Color(0), Background: termenv.ANSI256Color(22)}
	colorContentNew := Color{Foreground: termenv.ANSI256Color(0), Background: termenv.ANSI256Color(11)}
	strNewContent := "----- foo -----"

	t.Run("new content smaller than input", func(t *testing.T) {
		input := colorContent1.Sprint("Ghost")
		newContent := colorContentNew.Sprint("--")
		insertIdx := 2
		expectedOutput := "\x1b[38;5;0;48;5;12mGh\x1b[0m\x1b[38;5;0;48;5;11m--\x1b[0m\x1b[38;5;0;48;5;12mt\x1b[0m"
		output := overwriteContents(input, newContent, insertIdx, 80)
		assert.Equal(t, expectedOutput, output, fmt.Sprintf("newContent=%s", newContent))
	})

	t.Run("new content smaller than input with 2 colors", func(t *testing.T) {
		input := colorContent1.Sprint("Gho") + colorContent2.Sprint("st")
		newContent := colorContentNew.Sprint("--")
		insertIdx := 2
		expectedOutput := "\x1b[38;5;0;48;5;12mGh\x1b[0m\x1b[38;5;0;48;5;11m--\x1b[0m\x1b[38;5;0;48;5;22mt\x1b[0m"
		output := overwriteContents(input, newContent, insertIdx, 80)
		assert.Equal(t, expectedOutput, output, fmt.Sprintf("newContent=%s", newContent))
	})

	t.Run("new content longer than input", func(t *testing.T) {
		input := colorContent1.Sprint("Ghost")
		newContent := colorContentNew.Sprint(strNewContent)
		insertIdx := 2
		expectedOutput := "\x1b[38;5;0;48;5;12mGh\x1b[0m\x1b[38;5;0;48;5;11m----- foo -----\x1b[0m"
		output := overwriteContents(input, newContent, insertIdx, 80)
		assert.Equal(t, expectedOutput, output, fmt.Sprintf("newContent=%s", newContent))
	})

	t.Run("new content beyond input", func(t *testing.T) {
		input := colorContent1.Sprint("Ghost")
		newContent := colorContentNew.Sprint(strNewContent)
		insertIdx := 25
		expectedOutput := "\x1b[38;5;0;48;5;12mGhost\x1b[0m                    \x1b[38;5;0;48;5;11m----- foo -----\x1b[0m"
		output := overwriteContents(input, newContent, insertIdx, 80)
		assert.Equal(t, expectedOutput, output, fmt.Sprintf("newContent=%s", newContent))
	})

	t.Run("new content beyond input 2", func(t *testing.T) {
		input := "\x1b[38;5;237;48;5;233m 2 \x1b[0m   \x1b[0m\x1b[38;5;231m(\x1b[0m\x1b[38;5;186m'Arya'\x1b[0m\x1b[38;5;231m,\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;186m'Stark'\x1b[0m\x1b[38;5;231m,\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;141m3000\x1b[0m\x1b[38;5;231m),\x1b[0m"
		newContent := colorContentNew.Sprint("values")
		insertIdx := 58
		expectedOutput := "\x1b[38;5;237;48;5;233m 2 \x1b[0m   \x1b[0m\x1b[38;5;231m(\x1b[0m\x1b[38;5;186m'Arya'\x1b[0m\x1b[38;5;231m,\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;186m'Stark'\x1b[0m\x1b[38;5;231m,\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;141m3000\x1b[0m\x1b[38;5;231m),\x1b[0m                            \x1b[38;5;0;48;5;11mvalues\x1b[0m"
		output := overwriteContents(input, newContent, insertIdx, 80)
		assert.Equal(t, expectedOutput, output, fmt.Sprintf("newContent=%s", newContent))
	})

	t.Run("new content needs to be moved left", func(t *testing.T) {
		input := colorContent1.Sprint("Ghost")
		newContent := colorContentNew.Sprint(strNewContent)
		insertIdx := 2
		expectedOutput := "\x1b[38;5;0;48;5;12mG\x1b[0m\x1b[38;5;0;48;5;11m----- foo -----\x1b[0m"
		output := overwriteContents(input, newContent, insertIdx, 16)
		assert.Equal(t, expectedOutput, output, fmt.Sprintf("newContent=%s", newContent))
	})

	t.Run("new content needs to be moved left more", func(t *testing.T) {
		input := colorContent1.Sprint("Ghost")
		newContent := colorContentNew.Sprint(strNewContent)
		insertIdx := 2
		expectedOutput := "\x1b[38;5;0;48;5;11m----- foo -----\x1b[0m"
		output := overwriteContents(input, newContent, insertIdx, 15)
		assert.Equal(t, expectedOutput, output, fmt.Sprintf("newContent=%s", newContent))
	})

	t.Run("new content longer than display width", func(t *testing.T) {
		input := colorContent1.Sprint("Ghost")
		newContent := colorContentNew.Sprint(strNewContent)
		insertIdx := 2
		expectedOutput := "\x1b[38;5;0;48;5;11m----- foo \x1b[0m"
		output := overwriteContents(input, newContent, insertIdx, 10)
		assert.Equal(t, expectedOutput, output, fmt.Sprintf("newContent=%s", newContent))
	})
}

func Benchmark_stringSubset(b *testing.B) {
	color1 := Color{
		Foreground: termenv.ANSI256Color(0),
		Background: termenv.ANSI256Color(12),
	}

	input := color1.Sprint("Ghost")
	assert.Equal(b, "\x1b[38;5;0;48;5;12mGho\x1b[0m", stringSubset(input, 0, 2))
	for idx := 0; idx < b.N; idx++ {
		stringSubset(input, 0, 2)
	}
}

func Test_stringSubset(t *testing.T) {
	color1 := Color{
		Foreground: termenv.ANSI256Color(0),
		Background: termenv.ANSI256Color(12),
	}
	color2 := Color{
		Foreground: termenv.ANSI256Color(0),
		Background: termenv.ANSI256Color(22),
	}

	t.Run("single color/esc seq", func(t *testing.T) {
		input := color1.Sprint("Ghost")
		assert.Equal(t, "", stringSubset(input, 0, -1))
		assert.Equal(t, color1.Sprint("G"), stringSubset(input, 0, 0))
		assert.Equal(t, color1.Sprint("Gh"), stringSubset(input, 0, 1))
		assert.Equal(t, color1.Sprint("Gho"), stringSubset(input, 0, 2))
		assert.Equal(t, color1.Sprint("Ghos"), stringSubset(input, 0, 3))
		assert.Equal(t, color1.Sprint("Ghost"), stringSubset(input, 0, 4))
		assert.Equal(t, color1.Sprint("Ghost"), stringSubset(input, 0, 5))
		assert.Equal(t, color1.Sprint("host"), stringSubset(input, 1, 5))
		assert.Equal(t, color1.Sprint("ost"), stringSubset(input, 2, 5))
		assert.Equal(t, color1.Sprint("st"), stringSubset(input, 3, 5))
		assert.Equal(t, color1.Sprint("t"), stringSubset(input, 4, 5))
		assert.Equal(t, "", stringSubset(input, 5, 4))
	})

	t.Run("dual color/esc seq", func(t *testing.T) {
		input := color1.Sprint("Gho") + color2.Sprint("st")
		assert.Equal(t, "", stringSubset(input, 0, -1))
		assert.Equal(t, color1.Sprint("G"), stringSubset(input, 0, 0))
		assert.Equal(t, color1.Sprint("Gh"), stringSubset(input, 0, 1))
		assert.Equal(t, color1.Sprint("Gho"), stringSubset(input, 0, 2))
		assert.Equal(t, color1.Sprint("Gho")+color2.Sprint("s"), stringSubset(input, 0, 3))
		assert.Equal(t, color1.Sprint("Gho")+color2.Sprint("st"), stringSubset(input, 0, 4))
		assert.Equal(t, color1.Sprint("Gho")+color2.Sprint("st"), stringSubset(input, 0, 5))
		assert.Equal(t, color1.Sprint("ho")+color2.Sprint("st"), stringSubset(input, 1, 5))
		assert.Equal(t, color1.Sprint("o")+color2.Sprint("st"), stringSubset(input, 2, 5))
		assert.Equal(t, color2.Sprint("st"), stringSubset(input, 3, 5))
		assert.Equal(t, color2.Sprint("t"), stringSubset(input, 4, 5))
		assert.Equal(t, "", stringSubset(input, 5, 4))
	})
}
