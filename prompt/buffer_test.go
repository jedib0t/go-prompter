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

func TestBuffer_DeleteBackward(t *testing.T) {
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

func TestBuffer_DeleteForward(t *testing.T) {
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

func TestBuffer_DeleteWordBackward(t *testing.T) {
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

func TestBuffer_DeleteWordForward(t *testing.T) {
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

func TestBuffer_HasChanges(t *testing.T) {
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

func TestBuffer_Insert(t *testing.T) {
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
}

func TestBuffer_IsDone(t *testing.T) {
	b := getNewBuffer(t)
	assert.False(t, b.IsDone())

	b.done = true
	assert.True(t, b.IsDone())
}

func TestBuffer_Length(t *testing.T) {
	b := getNewBuffer(t)
	assert.Equal(t, 0, b.Length())

	b.lines = []string{"abc"}
	assert.Equal(t, 3, b.Length())

	b.lines = []string{"abc", ""}
	assert.Equal(t, 4, b.Length())

	b.lines = []string{"abc", "def"}
	assert.Equal(t, 7, b.Length())
}

func TestBuffer_MakeWordCapitalCase(t *testing.T) {
	b := getNewBuffer(t)
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

func TestBuffer_MakeWordLowerCase(t *testing.T) {
	b := getNewBuffer(t)
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

func TestBuffer_MakeWordUpperCase(t *testing.T) {
	b := getNewBuffer(t)
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

func TestBuffer_MarkAsDone(t *testing.T) {
	b := getNewBuffer(t)
	assert.False(t, b.done)

	b.MarkAsDone()
	assert.True(t, b.done)
}

func TestBuffer_MoveDown(t *testing.T) {
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

func TestBuffer_MoveLeft(t *testing.T) {
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

func TestBuffer_MoveLineBegin(t *testing.T) {
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

func TestBuffer_MoveLineEnd(t *testing.T) {
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

func TestBuffer_MoveRight(t *testing.T) {
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

func TestBuffer_MoveUp(t *testing.T) {
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

func TestBuffer_MoveWordLeft(t *testing.T) {
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

func TestBuffer_MoveWordRight(t *testing.T) {
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

func TestBuffer_Render(t *testing.T) {

}

func TestBuffer_Set(t *testing.T) {
	b := getNewBuffer(t)
	b.tab = "    "

	b.Set("echo $VARIABLE\necho $VARIABLE2\t#testing")
	assert.Equal(t, []string{"echo $VARIABLE", "echo $VARIABLE2    #testing"}, b.lines)
	assert.Equal(t, CursorLocation{Line: 1, Column: 27}, b.cursor)
}

func TestBuffer_String(t *testing.T) {
	b := getNewBuffer(t)
	assert.Equal(t, "", b.String())

	b.lines = []string{"abc"}
	assert.Equal(t, "abc", b.String())

	b.lines = []string{"abc", ""}
	assert.Equal(t, "abc\n", b.String())

	b.lines = []string{"abc", "def"}
	assert.Equal(t, "abc\ndef", b.String())
}

func TestBuffer_SwapCharacterNext(t *testing.T) {

}

func TestBuffer_SwapCharacterPrevious(t *testing.T) {

}

func TestBuffer_SwapWordNext(t *testing.T) {

}

func TestBuffer_SwapWordPrevious(t *testing.T) {

}
