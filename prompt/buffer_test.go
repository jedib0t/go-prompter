package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getNewBuffer(t *testing.T) *buffer {
	b := newBuffer()
	assert.Equal(t, []string{""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
	return b
}

func TestBufferCursor(t *testing.T) {
	b := getNewBuffer(t)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.Cursor())

	b.Set("foo\nbar")
	assert.Equal(t, CursorLocation{Line: 1, Column: 3}, b.Cursor())
}

func TestBufferDeleteBackward(t *testing.T) {
	b := getNewBuffer(t)

	b.lines = []string{"abc", "def"}
	b.cursor = CursorLocation{Line: 1, Column: 1}
	b.DeleteBackward(1)
	assert.Equal(t, []string{"abc", "ef"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)
	b.DeleteBackward(1)
	assert.Equal(t, []string{"abcef"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 3}, b.cursor)
	b.DeleteBackward(1)
	assert.Equal(t, []string{"abef"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)
	b.DeleteBackward(1)
	assert.Equal(t, []string{"aef"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 1}, b.cursor)
	b.DeleteBackward(1)
	assert.Equal(t, []string{"ef"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
	b.DeleteBackward(1)
	assert.Equal(t, []string{"ef"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 2}
	b.DeleteBackward(1)
	assert.Equal(t, []string{"e"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 1}, b.cursor)
	b.DeleteBackward(1)
	assert.Equal(t, []string{""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
	b.DeleteBackward(1)
	assert.Equal(t, []string{""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.lines = []string{"abc", "def", "ghi"}
	b.cursor = CursorLocation{Line: 1, Column: 1}
	b.DeleteBackward(-1)
	assert.Equal(t, []string{"ef", "ghi"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.lines = []string{"abc", "def", "ghi"}
	b.cursor = CursorLocation{Line: 2, Column: 3}
	b.DeleteBackward(-1)
	assert.Equal(t, []string{""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
}

func TestBufferDeleteBackwardToBeginningOfLine(t *testing.T) {
	b := getNewBuffer(t)
	b.InsertString("foo bar baz")
	b.MoveWordLeft()
	b.DeleteBackwardToBeginningOfLine()

	assert.Equal(t, "baz", b.String())
}

func TestBufferDeleteForward(t *testing.T) {
	b := getNewBuffer(t)

	b.lines = []string{"abc", "def"}
	b.cursor = CursorLocation{Line: 0, Column: 2}
	b.DeleteForward(1)
	assert.Equal(t, []string{"ab", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)
	b.DeleteForward(1)
	assert.Equal(t, []string{"abdef"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)
	b.DeleteForward(1)
	assert.Equal(t, []string{"abef"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)
	b.DeleteForward(1)
	assert.Equal(t, []string{"abf"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)
	b.DeleteForward(1)
	assert.Equal(t, []string{"ab"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)
	b.DeleteForward(1)
	assert.Equal(t, []string{"ab"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.DeleteForward(1)
	assert.Equal(t, []string{"b"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
	b.DeleteForward(1)
	assert.Equal(t, []string{""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
	b.DeleteForward(1)
	assert.Equal(t, []string{""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.lines = []string{"abc", "def", "ghi"}
	b.cursor = CursorLocation{Line: 0, Column: 3}
	b.DeleteForward(1)
	assert.Equal(t, []string{"abcdef", "ghi"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 3}, b.cursor)

	b.lines = []string{"abc", "def"}
	b.cursor = CursorLocation{Line: 0, Column: 3}
	b.DeleteForward(-1)
	assert.Equal(t, []string{"abc"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 3}, b.cursor)

	b.lines = []string{"abc", "def", "ghi"}
	b.cursor = CursorLocation{Line: 0, Column: 2}
	b.DeleteForward(-1)
	assert.Equal(t, []string{"ab"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)
}

func TestBufferDeleteForwardToEndOfLine(t *testing.T) {
	b := getNewBuffer(t)
	b.InsertString("foo baz bar")
	b.MoveWordLeft()
	b.MoveWordLeft()
	b.DeleteForwardToEndOfLine()

	assert.Equal(t, "foo ", b.String())
}

func TestBufferDeleteWordBackward(t *testing.T) {
	b := getNewBuffer(t)

	b.lines = []string{"abc def"}
	b.cursor = CursorLocation{Line: 0, Column: 4}
	b.DeleteWordBackward()
	assert.Equal(t, []string{"def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
	b.DeleteWordBackward()
	assert.Equal(t, []string{"def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.lines = []string{"abc def"}
	b.cursor = CursorLocation{Line: 0, Column: 3}
	b.DeleteWordBackward()
	assert.Equal(t, []string{" def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
	b.DeleteWordBackward()
	assert.Equal(t, []string{" def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.lines = []string{"abc def ghi"}
	b.cursor = CursorLocation{Line: 0, Column: 11}
	b.DeleteWordBackward()
	assert.Equal(t, []string{"abc def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 7}, b.cursor)
	b.DeleteWordBackward()
	assert.Equal(t, []string{"abc"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 3}, b.cursor)
	b.DeleteWordBackward()
	assert.Equal(t, []string{""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.lines = []string{"abc", "def"}
	b.cursor = CursorLocation{Line: 1, Column: 3}
	b.DeleteWordBackward()
	assert.Equal(t, []string{"abc", ""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)
	b.DeleteWordBackward()
	assert.Equal(t, []string{"abc"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 3}, b.cursor)

	b.lines = []string{"abc", "def", "ghi"}
	b.cursor = CursorLocation{Line: 2, Column: 3}
	b.DeleteWordBackward()
	assert.Equal(t, []string{"abc", "def", ""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 2, Column: 0}, b.cursor)
	b.DeleteWordBackward()
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 3}, b.cursor)
	b.DeleteWordBackward()
	assert.Equal(t, []string{"abc", ""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)
	b.DeleteWordBackward()
	assert.Equal(t, []string{"abc"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 3}, b.cursor)
	b.DeleteWordBackward()
	assert.Equal(t, []string{""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
	b.DeleteWordBackward()
	assert.Equal(t, []string{""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
}

func TestBufferDeleteWordForward(t *testing.T) {
	b := getNewBuffer(t)

	b.lines = []string{"abc def"}
	b.cursor = CursorLocation{Line: 0, Column: 4}
	b.DeleteWordForward()
	assert.Equal(t, []string{"abc "}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 4}, b.cursor)
	b.DeleteWordForward()
	assert.Equal(t, []string{"abc "}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 4}, b.cursor)

	b.lines = []string{"abc def"}
	b.cursor = CursorLocation{Line: 0, Column: 2}
	b.DeleteWordForward()
	assert.Equal(t, []string{"abdef"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)
	b.DeleteWordForward()
	assert.Equal(t, []string{"ab"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)
	b.DeleteWordForward()
	assert.Equal(t, []string{"ab"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)

	b.lines = []string{"abc", "def"}
	b.cursor = CursorLocation{Line: 0, Column: 2}
	b.DeleteWordForward()
	assert.Equal(t, []string{"ab", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)
	b.DeleteWordForward()
	assert.Equal(t, []string{"abdef"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)

	b.lines = []string{"abc", "def", "ghi"}
	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.DeleteWordForward()
	assert.Equal(t, []string{"", "def", "ghi"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
	b.DeleteWordForward()
	assert.Equal(t, []string{"def", "ghi"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
	b.DeleteWordForward()
	assert.Equal(t, []string{"", "ghi"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
	b.DeleteWordForward()
	assert.Equal(t, []string{"ghi"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
	b.DeleteWordForward()
	assert.Equal(t, []string{""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
	b.DeleteWordForward()
	assert.Equal(t, []string{""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
}

func TestBufferDisplay(t *testing.T) {
	b := getNewBuffer(t)

	lines, cur := b.Display()
	assert.Equal(t, []string{""}, lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, cur)

	b.InsertString("bar\nbaz")
	lines, cur = b.Display()
	assert.Equal(t, []string{"bar", "baz"}, lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 3}, cur)
}

func TestBufferHasChanges(t *testing.T) {
	b := getNewBuffer(t)

	b.linesChanged.Clear()
	assert.False(t, b.HasChanges())

	b.InsertString("abc")
	assert.True(t, b.HasChanges())

	b.linesChanged.Clear()
	assert.False(t, b.HasChanges())

	b.MoveToBeginning()
	assert.True(t, b.HasChanges())

	b.linesChanged.Clear()
	assert.False(t, b.HasChanges())
}

func TestBufferInsert(t *testing.T) {
	b := getNewBuffer(t)

	b.Insert('\n')
	assert.Equal(t, []string{"", ""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)

	b.lines = []string{""}
	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.Insert('a')
	assert.Equal(t, []string{"a"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 1}, b.cursor)

	b.Insert('b')
	assert.Equal(t, []string{"ab"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)

	b.Insert('c')
	assert.Equal(t, []string{"abc"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 3}, b.cursor)

	b.Insert('\n')
	assert.Equal(t, []string{"abc", ""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)

	b.InsertString("    ") // tab gets inserted as "    "
	assert.Equal(t, []string{"abc", "    "}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 4}, b.cursor)

	b.Insert('\n')
	assert.Equal(t, []string{"abc", "    ", ""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 2, Column: 0}, b.cursor)

	b.Insert('d')
	assert.Equal(t, []string{"abc", "    ", "d"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 2, Column: 1}, b.cursor)

	b.Insert('e')
	assert.Equal(t, []string{"abc", "    ", "de"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 2, Column: 2}, b.cursor)

	b.Insert('f')
	assert.Equal(t, []string{"abc", "    ", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 2, Column: 3}, b.cursor)

	b.cursor = CursorLocation{Line: 1, Column: 2}
	b.Insert('\n')
	assert.Equal(t, []string{"abc", "  ", "  ", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 2, Column: 0}, b.cursor)

	b.Insert('1')
	assert.Equal(t, []string{"abc", "  ", "1  ", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 2, Column: 1}, b.cursor)

	b.cursor = CursorLocation{Line: 2, Column: 3}
	b.Insert('2')
	assert.Equal(t, []string{"abc", "  ", "1  2", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 2, Column: 4}, b.cursor)

	b.Insert('\n')
	assert.Equal(t, []string{"abc", "  ", "1  2", "", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 3, Column: 0}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.Insert('?')
	assert.Equal(t, []string{"?abc", "  ", "1  2", "", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 1}, b.cursor)

	b.SetTab("    ")
	b.Insert('\t')
	assert.Equal(t, []string{"?    abc", "  ", "1  2", "", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 5}, b.cursor)
}

func TestBufferIsDone(t *testing.T) {
	b := getNewBuffer(t)
	assert.False(t, b.IsDone())

	b.done = true
	assert.True(t, b.IsDone())
}

func TestBufferLength(t *testing.T) {
	b := getNewBuffer(t)
	assert.Equal(t, 0, b.Length())

	b.lines = []string{"abc"}
	assert.Equal(t, 3, b.Length())

	b.lines = []string{"abc", ""}
	assert.Equal(t, 4, b.Length())

	b.lines = []string{"abc", "def"}
	assert.Equal(t, 7, b.Length())
}

func TestBufferLines(t *testing.T) {
	b := getNewBuffer(t)
	assert.Equal(t, []string{""}, b.Lines())

	b.InsertString("foo\n")
	assert.Equal(t, []string{"foo", ""}, b.Lines())

	b.InsertString("bar\n")
	assert.Equal(t, []string{"foo", "bar", ""}, b.Lines())
}

func TestBufferMakeWordCapitalCase(t *testing.T) {
	b := getNewBuffer(t)
	b.MakeWordCapitalCase()

	b.lines = []string{"abc", "def", "ghi"}
	b.cursor = CursorLocation{Line: 0, Column: 0}

	b.MakeWordCapitalCase()
	assert.Equal(t, []string{"Abc", "def", "ghi"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)
	b.cursor.Column = 1
	b.MakeWordCapitalCase()
	assert.Equal(t, []string{"Abc", "Def", "ghi"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 2, Column: 0}, b.cursor)
	b.cursor.Column = 2
	b.MakeWordCapitalCase()
	assert.Equal(t, []string{"Abc", "Def", "Ghi"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 2, Column: 3}, b.cursor)
}

func TestBufferMakeWordLowerCase(t *testing.T) {
	b := getNewBuffer(t)
	b.MakeWordLowerCase()
	b.lines = []string{"ABC", "DEF", "GHI"}
	b.cursor = CursorLocation{Line: 0, Column: 0}

	b.MakeWordLowerCase()
	assert.Equal(t, []string{"abc", "DEF", "GHI"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)
	b.cursor.Column = 1
	b.MakeWordLowerCase()
	assert.Equal(t, []string{"abc", "def", "GHI"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 2, Column: 0}, b.cursor)
	b.cursor.Column = 2
	b.MakeWordLowerCase()
	assert.Equal(t, []string{"abc", "def", "ghi"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 2, Column: 3}, b.cursor)
}

func TestBufferMakeWordUpperCase(t *testing.T) {
	b := getNewBuffer(t)
	b.MakeWordUpperCase()
	b.lines = []string{"abc", "def", "ghi"}
	b.cursor = CursorLocation{Line: 0, Column: 0}

	b.MakeWordUpperCase()
	assert.Equal(t, []string{"ABC", "def", "ghi"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)
	b.cursor.Column = 1
	b.MakeWordUpperCase()
	assert.Equal(t, []string{"ABC", "DEF", "ghi"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 2, Column: 0}, b.cursor)
	b.cursor.Column = 2
	b.MakeWordUpperCase()
	assert.Equal(t, []string{"ABC", "DEF", "GHI"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 2, Column: 3}, b.cursor)
}

func TestBufferMarkAsDone(t *testing.T) {
	b := getNewBuffer(t)
	assert.False(t, b.done)

	b.MarkAsDone()
	assert.True(t, b.done)
}

func TestBufferMoveDown(t *testing.T) {
	b := getNewBuffer(t)

	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.MoveDown(1)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.lines = []string{"abc123", "def"}
	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.MoveDown(1)
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 2}
	b.MoveDown(1)
	assert.Equal(t, CursorLocation{Line: 1, Column: 2}, b.cursor)
	b.MoveDown(1)
	assert.Equal(t, CursorLocation{Line: 1, Column: 2}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 5}
	b.MoveDown(1)
	assert.Equal(t, CursorLocation{Line: 1, Column: 3}, b.cursor)
	b.MoveDown(1)
	assert.Equal(t, CursorLocation{Line: 1, Column: 3}, b.cursor)
}

func TestBufferMoveLeft(t *testing.T) {
	b := getNewBuffer(t)

	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.MoveLeft(1)
	assert.Equal(t, []string{""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.lines = []string{"abc", "def"}
	b.cursor = CursorLocation{Line: 1, Column: 3}
	b.MoveLeft(-1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.cursor = CursorLocation{Line: 1, Column: 3}
	b.MoveLeft(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 2}, b.cursor)

	b.MoveLeft(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 1}, b.cursor)

	b.MoveLeft(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)

	b.MoveLeft(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 3}, b.cursor)

	b.MoveLeft(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)

	b.MoveLeft(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 1}, b.cursor)

	b.MoveLeft(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.MoveLeft(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
}

func TestBufferMoveRight(t *testing.T) {
	b := getNewBuffer(t)

	b.MoveRight(1)
	assert.Equal(t, []string{""}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.lines = []string{"abc", "def"}
	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.MoveRight(-1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 3}, b.cursor)

	b.cursor = CursorLocation{Line: 1, Column: 3}
	b.MoveRight(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 3}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.MoveRight(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 1}, b.cursor)

	b.MoveRight(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)

	b.MoveRight(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 0, Column: 3}, b.cursor)

	b.MoveRight(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)

	b.MoveRight(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 1}, b.cursor)

	b.MoveRight(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 2}, b.cursor)

	b.MoveRight(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 3}, b.cursor)

	b.MoveRight(1)
	assert.Equal(t, []string{"abc", "def"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 3}, b.cursor)
}

func TestBufferMoveToBeginning(t *testing.T) {
	b := getNewBuffer(t)

	b.InsertString("foo")
	b.MoveToBeginning()
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.Cursor())
}

func TestBufferMoveToBeginningOfLine(t *testing.T) {
	b := getNewBuffer(t)

	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.MoveToBeginningOfLine()
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.lines = []string{"abc", "def"}
	b.cursor = CursorLocation{Line: 1, Column: 3}
	b.MoveToBeginningOfLine()
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)

	b.cursor = CursorLocation{Line: 1, Column: 2}
	b.MoveToBeginningOfLine()
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 3}
	b.MoveToBeginningOfLine()
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
}

func TestBufferMoveToEnd(t *testing.T) {
	b := getNewBuffer(t)

	b.InsertString("foo")
	b.MoveToEnd()
	assert.Equal(t, CursorLocation{Line: 0, Column: 3}, b.Cursor())
}

func TestBufferMoveToEndOfLine(t *testing.T) {
	b := getNewBuffer(t)

	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.MoveToEndOfLine()
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.lines = []string{"abc", "def"}
	b.cursor = CursorLocation{Line: 1, Column: 3}
	b.MoveToEndOfLine()
	assert.Equal(t, CursorLocation{Line: 1, Column: 3}, b.cursor)

	b.cursor = CursorLocation{Line: 1, Column: 2}
	b.MoveToEndOfLine()
	assert.Equal(t, CursorLocation{Line: 1, Column: 3}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 3}
	b.MoveToEndOfLine()
	assert.Equal(t, CursorLocation{Line: 0, Column: 3}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 2}
	b.MoveToEndOfLine()
	assert.Equal(t, CursorLocation{Line: 0, Column: 3}, b.cursor)
}

func TestBufferMoveUp(t *testing.T) {
	b := getNewBuffer(t)

	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.MoveUp(1)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.lines = []string{"abc", "def123"}
	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.MoveUp(1)
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.cursor = CursorLocation{Line: 1, Column: 2}
	b.MoveUp(1)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)
	b.MoveUp(1)
	assert.Equal(t, CursorLocation{Line: 0, Column: 2}, b.cursor)

	b.cursor = CursorLocation{Line: 1, Column: 5}
	b.MoveUp(1)
	assert.Equal(t, CursorLocation{Line: 0, Column: 3}, b.cursor)
	b.MoveUp(1)
	assert.Equal(t, CursorLocation{Line: 0, Column: 3}, b.cursor)
}

func TestBufferMoveWordLeft(t *testing.T) {
	b := getNewBuffer(t)

	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.MoveWordLeft()
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.lines = []string{"abc 123", "def"}
	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.MoveWordLeft()
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 1}
	b.MoveWordLeft()
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 2}
	b.MoveWordLeft()
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 3}
	b.MoveWordLeft()
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 4}
	b.MoveWordLeft()
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.cursor = CursorLocation{Line: 1, Column: 3}
	b.MoveWordLeft()
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)
	b.MoveWordLeft()
	assert.Equal(t, CursorLocation{Line: 0, Column: 4}, b.cursor)
	b.MoveWordLeft()
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
}

func TestBufferMoveWordRight(t *testing.T) {
	b := getNewBuffer(t)

	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.MoveWordRight()
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)

	b.lines = []string{"abc 123 ", "def"}
	b.cursor = CursorLocation{Line: 0, Column: 0}
	b.MoveWordRight()
	assert.Equal(t, CursorLocation{Line: 0, Column: 4}, b.cursor)
	b.MoveWordRight()
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)
	b.MoveWordRight()
	assert.Equal(t, CursorLocation{Line: 1, Column: 3}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 1}
	b.MoveWordRight()
	assert.Equal(t, CursorLocation{Line: 0, Column: 4}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 2}
	b.MoveWordRight()
	assert.Equal(t, CursorLocation{Line: 0, Column: 4}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 3}
	b.MoveWordRight()
	assert.Equal(t, CursorLocation{Line: 0, Column: 4}, b.cursor)

	b.cursor = CursorLocation{Line: 0, Column: 8}
	b.MoveWordRight()
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, b.cursor)
}

func TestBufferNumLines(t *testing.T) {
	b := getNewBuffer(t)
	assert.Equal(t, 1, b.NumLines())

	b.InsertString("food\nbard")
	assert.Equal(t, 2, b.NumLines())
}

func TestBufferReset(t *testing.T) {
	b := getNewBuffer(t)
	b.InsertString("foo")

	b.Reset()
	assert.Len(t, b.lines, 1)
	assert.Empty(t, b.lines[0])
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, b.cursor)
}

func TestBufferSet(t *testing.T) {
	b := getNewBuffer(t)
	b.tab = "    "

	b.Set("echo $VARIABLE\necho $VARIABLE2\t#testing")
	assert.Equal(t, []string{"echo $VARIABLE", "echo $VARIABLE2    #testing"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 27}, b.cursor)
}

func TestBufferString(t *testing.T) {
	b := getNewBuffer(t)
	assert.Equal(t, "", b.String())

	b.lines = []string{"abc"}
	assert.Equal(t, "abc", b.String())

	b.lines = []string{"abc", ""}
	assert.Equal(t, "abc\n", b.String())

	b.lines = []string{"abc", "def"}
	assert.Equal(t, "abc\ndef", b.String())
}

func TestBufferGetWordAtCursor(t *testing.T) {
	b := getNewBuffer(t)
	b.InsertString("foo bar baz foo")

	b.cursor.Column = 11
	word, idx := b.getWordAtCursor()
	assert.Equal(t, "baz", word)
	assert.Equal(t, 8, idx)

	b.cursor.Column = 7
	word, idx = b.getWordAtCursor()
	assert.Equal(t, "bar", word)
	assert.Equal(t, 4, idx)

	b.cursor.Column = 3
	word, idx = b.getWordAtCursor()
	assert.Equal(t, "foo", word)
	assert.Equal(t, 0, idx)

	b.cursor.Column = 2
	word, idx = b.getWordAtCursor()
	assert.Equal(t, "", word)
	assert.Equal(t, -1, idx)
}

func TestLinesChangedMap(t *testing.T) {
	lcm := make(linesChangedMap)
	assert.Empty(t, lcm)
	assert.True(t, lcm.NothingChanged())

	lcm.Mark(1)
	assert.False(t, lcm.AllChanged())
	assert.True(t, lcm.AnythingChanged())
	assert.False(t, lcm.IsChanged(0))
	assert.True(t, lcm.IsChanged(1))
	assert.False(t, lcm.IsChanged(2))
	assert.False(t, lcm.NothingChanged())

	lcm.MarkAll()
	assert.True(t, lcm.AllChanged())
	assert.True(t, lcm.IsChanged(2))
	assert.False(t, lcm.NothingChanged())

	assert.Equal(t, "[-1 1]", lcm.String())
}
