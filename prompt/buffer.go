package prompt

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// buffer helps store the user input, track the cursor position, and help
// manipulate the user input with adding/removing strings from it
type buffer struct {
	cursor         CursorLocation
	cursorRendered CursorLocation
	done           bool
	lines          []string
	linesChanged   linesChangedMap
	linesRendered  string
	mutex          sync.Mutex
	tab            string
}

// newBuffer returns a buffer object with sane defaults
func newBuffer() *buffer {
	return &buffer{
		cursor:        CursorLocation{Line: 0, Column: 0},
		done:          false,
		lines:         []string{""},
		linesChanged:  make(linesChangedMap),
		linesRendered: fmt.Sprint(time.Now().Format(time.RFC3339Nano)),
		mutex:         sync.Mutex{},
	}
}

// Cursor returns the current Cursor Location.
func (b *buffer) Cursor() CursorLocation {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.cursor
}

// DeleteBackward deletes n runes backwards
func (b *buffer) DeleteBackward(n int, locked ...bool) {
	if len(locked) == 0 {
		b.mutex.Lock()
		defer b.mutex.Unlock()
	}

	// if asked to delete till beginning, just set N to the max value possible
	if n == -1 {
		n = len(strings.Join(b.lines, "\n"))
	}

	// delete backward rune by rune
	for ; n > 0; n-- {
		if b.cursor.Column == 0 {
			if b.cursor.Line > 0 {
				prevLine, line := b.getLine(b.cursor.Line-1), b.getLine(b.cursor.Line)
				var lines []string
				lines = append(lines, b.lines[:b.cursor.Line-1]...)
				lines = append(lines, prevLine+line)
				if b.cursor.Line < len(b.lines)-1 {
					lines = append(lines, b.lines[b.cursor.Line+1:]...)
				}

				b.lines = lines
				b.linesChanged.MarkAll()
				b.cursor.Line--
				b.cursor.Column = len(prevLine)
			}
		} else {
			line := b.getCurrentLine()
			b.lines[b.cursor.Line] = line[:b.cursor.Column-1] + line[b.cursor.Column:]
			b.linesChanged.Mark(b.cursor.Line)
			b.cursor.Column--
		}
	}
}

// DeleteBackwardToBeginningOfLine deletes till cursor reaches 0th column
func (b *buffer) DeleteBackwardToBeginningOfLine(locked ...bool) {
	if len(locked) == 0 {
		b.mutex.Lock()
		defer b.mutex.Unlock()
	}

	b.DeleteBackward(b.cursor.Column, true)
}

// DeleteForward deletes n runes forwards
func (b *buffer) DeleteForward(n int, locked ...bool) {
	if len(locked) == 0 {
		b.mutex.Lock()
		defer b.mutex.Unlock()
	}

	// if asked to delete till end, just set N to the max value possible
	if n == -1 {
		n = len(strings.Join(b.lines, "\n"))
	}

	// delete forward rune by rune
	for ; n > 0; n-- {
		line := b.getCurrentLine()
		if b.cursor.Column == len(line) {
			if b.cursor.Line == len(b.lines)-1 {
				return
			}
			line += b.getLine(b.cursor.Line + 1)

			var lines []string
			lines = append(lines, b.lines[:b.cursor.Line]...)
			lines = append(lines, line)
			if b.cursor.Line < len(b.lines)-2 {
				lines = append(lines, b.lines[b.cursor.Line+2:]...)
			}

			b.lines = lines
			b.linesChanged.MarkAll()
		} else if b.cursor.Column > 0 {
			b.lines[b.cursor.Line] = line[:b.cursor.Column] + line[b.cursor.Column+1:]
			b.linesChanged.Mark(b.cursor.Line)
		} else {
			b.lines[b.cursor.Line] = line[b.cursor.Column+1:]
			b.linesChanged.Mark(b.cursor.Line)
		}
	}
}

// DeleteForwardToEndOfLine deletes till cursor reaches 0th column
func (b *buffer) DeleteForwardToEndOfLine(locked ...bool) {
	if len(locked) == 0 {
		b.mutex.Lock()
		defer b.mutex.Unlock()
	}

	b.DeleteForward(len(b.getCurrentLine())-b.cursor.Column, true)
}

// DeleteWordBackward deletes the previous word
func (b *buffer) DeleteWordBackward() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// already on the first column? delete one char backwards...
	if b.cursor.Column == 0 {
		b.DeleteBackward(1, true)
		return
	}

	// delete till beginning of previous word
	foundWord := false
	line := b.getCurrentLine()
	for idx := b.cursor.Column - 1; idx >= 0; idx-- {
		isPartOfWord := isPartOfWord(line[idx])
		if !isPartOfWord && foundWord {
			b.lines[b.cursor.Line] = line[:idx] + line[b.cursor.Column:]
			b.cursor.Column = idx
			return
		}
		if isPartOfWord {
			foundWord = true
		}
	}
	b.lines[b.cursor.Line] = line[b.cursor.Column:]
	b.linesChanged.Mark(b.cursor.Line)
	b.cursor.Column = 0
}

// DeleteWordForward deletes the next word
func (b *buffer) DeleteWordForward() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// already on the last column? delete one char forwards...
	line := b.getCurrentLine()
	if b.cursor.Column == len(line) {
		b.DeleteForward(1, true)
		return
	}

	// delete till beginning of previous word
	foundWord, foundNonWord := false, false
	for idx := b.cursor.Column; idx < len(line); idx++ {
		isPartOfWord := isPartOfWord(line[idx])
		if !isPartOfWord {
			foundNonWord = true
		}
		if isPartOfWord && foundWord && foundNonWord {
			b.lines[b.cursor.Line] = line[:b.cursor.Column] + line[idx:]
			return
		}
		if isPartOfWord {
			foundWord = true
		}
	}
	b.lines[b.cursor.Line] = line[:b.cursor.Column]
	b.linesChanged.Mark(b.cursor.Line)
}

// Display returns the current contents of the buffer for display and assumes
// that all returned content has been displayed/rendered.
func (b *buffer) Display() ([]string, CursorLocation) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	lines := append([]string{}, b.lines...)
	cursor := b.cursor
	b.linesChanged.Clear()

	return lines, cursor
}

// HasChanges returns true if Render() will return something else on the next
// call to it.
func (b *buffer) HasChanges() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.linesChanged.AnythingChanged()
}

// Insert inserts the string at the current cursor position
func (b *buffer) Insert(r rune, locked ...bool) {
	if len(locked) == 0 {
		b.mutex.Lock()
		defer b.mutex.Unlock()
	}

	if r == '\n' {
		line := b.getCurrentLine()

		var lines []string
		lines = append(lines, b.lines[:b.cursor.Line]...)
		lines = append(lines, line[:b.cursor.Column])
		if b.cursor.Column < len(line) { // cursor somewhere before the end
			lines = append(lines, line[b.cursor.Column:])
		} else {
			lines = append(lines, "")
		}
		lines = append(lines, b.lines[b.cursor.Line+1:]...)

		b.lines = lines
		b.linesChanged.MarkAll()
		b.cursor.Line++
		b.cursor.Column = 0
	} else {
		rStr := fmt.Sprintf("%c", r)
		if r == '\t' {
			rStr = b.tab
		}

		line := b.getCurrentLine()
		line = line[:b.cursor.Column] + rStr + line[b.cursor.Column:]

		b.lines[b.cursor.Line] = line
		b.linesChanged.Mark(b.cursor.Line)
		b.cursor.Column += len(rStr)
	}
}

func (b *buffer) InsertString(str string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, r := range str {
		b.Insert(r, true)
	}
}

func (b *buffer) IsDone() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.done
}

// Length returns the current input length
func (b *buffer) Length() int {
	return len(b.String())
}

// Lines returns a copy of the lines in the buffer.
func (b *buffer) Lines() []string {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return append([]string{}, b.lines...)
}

// MakeWordCapitalCase converts the current word to Capital case
func (b *buffer) MakeWordCapitalCase() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	line := b.getCurrentLine()
	word, idxWordStart, idxWordEnd := b.getCurrentWord(line)
	if word == "" || idxWordStart == -1 || idxWordEnd == -1 {
		return
	}

	word = strings.ToUpper(word[0:1]) + word[1:]
	b.lines[b.cursor.Line] = line[:idxWordStart] + word + line[idxWordEnd:]
	b.linesChanged.Mark(b.cursor.Line)
	b.MoveWordRight(true)
}

// MakeWordLowerCase converts the current word to Lower case
func (b *buffer) MakeWordLowerCase() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	line := b.getCurrentLine()
	word, idxWordStart, idxWordEnd := b.getCurrentWord(line)
	if word == "" || idxWordStart == -1 || idxWordEnd == -1 {
		return
	}

	b.lines[b.cursor.Line] = line[:idxWordStart] + strings.ToLower(word) + line[idxWordEnd:]
	b.linesChanged.Mark(b.cursor.Line)
	b.MoveWordRight(true)
}

// MakeWordUpperCase converts the current word to Upper case
func (b *buffer) MakeWordUpperCase() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	line := b.getCurrentLine()
	word, idxWordStart, idxWordEnd := b.getCurrentWord(line)
	if word == "" || idxWordStart == -1 || idxWordEnd == -1 {
		return
	}

	b.lines[b.cursor.Line] = line[:idxWordStart] + strings.ToUpper(word) + line[idxWordEnd:]
	b.linesChanged.Mark(b.cursor.Line)
	b.MoveWordRight(true)
}

// MarkAsDone signifies that the user input is done
func (b *buffer) MarkAsDone() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.done = true
}

// MoveDown attempts to move the cursor to the same position in the next line
func (b *buffer) MoveDown(n int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.linesChanged.Mark(b.cursor.Line)
	b.cursor.Line += n
	if b.cursor.Line >= len(b.lines) {
		b.cursor.Line = len(b.lines) - 1
	}
	b.linesChanged.Mark(b.cursor.Line)
	line := b.getCurrentLine()
	if b.cursor.Column > len(line) {
		b.cursor.Column = len(line)
	}
}

// MoveLeft moves the cursor left n runes
func (b *buffer) MoveLeft(n int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// move to the very beginning
	if n == -1 {
		b.linesChanged.MarkAll()
		b.cursor = CursorLocation{Line: 0, Column: 0}
		return
	}

	// move left until n becomes 0, or beginning of buffer is reached
	for ; n > 0; n-- {
		b.cursor.Column--
		b.linesChanged.Mark(b.cursor.Line)
		if b.cursor.Column < 0 {
			if b.cursor.Line == 0 {
				b.cursor.Column = 0
				break
			}
			b.cursor.Line--
			b.linesChanged.Mark(b.cursor.Line)
			b.cursor.Column = len(b.getCurrentLine())
		}
	}
}

// MoveRight moves the cursor right n runes
func (b *buffer) MoveRight(n int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// move to the very end
	if n == -1 {
		b.cursor.Line = len(b.lines) - 1
		b.cursor.Column = len(b.getCurrentLine())
		return
	}

	// move right until n becomes 0, or end of buffer is reached
	for ; n > 0; n-- {
		line := b.getCurrentLine()
		b.cursor.Column++
		b.linesChanged.Mark(b.cursor.Line)
		if b.cursor.Column > len(line) {
			if b.cursor.Line == len(b.lines)-1 {
				b.cursor.Column = len(line)
				break
			}
			b.cursor.Line++
			b.linesChanged.Mark(b.cursor.Line)
			b.cursor.Column = 0
		}
	}
}

// MoveToBeginning moves the cursor right to the beginning of the first line
func (b *buffer) MoveToBeginning() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.linesChanged.Mark(b.cursor.Line) // before
	b.cursor.Line = 0
	b.linesChanged.Mark(b.cursor.Line) // after
	b.cursor.Column = 0
}

// MoveToBeginningOfLine moves the cursor right to the beginning of the current line
func (b *buffer) MoveToBeginningOfLine() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.cursor.Column = 0
	b.linesChanged.Mark(b.cursor.Line) // before
}

// MoveToEnd moves the cursor right to the end of the last line
func (b *buffer) MoveToEnd() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.linesChanged.Mark(b.cursor.Line) // before
	b.cursor.Line = len(b.lines) - 1
	b.linesChanged.Mark(b.cursor.Line) // after
	b.cursor.Column = len(b.getCurrentLine())
}

// MoveToEndOfLine moves the cursor right to the end of the current line
func (b *buffer) MoveToEndOfLine() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.cursor.Column = len(b.getCurrentLine())
	b.linesChanged.Mark(b.cursor.Line) // before
}

// MoveUp attempts to move the cursor to the same position in the previous line
func (b *buffer) MoveUp(n int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.linesChanged.Mark(b.cursor.Line) // before
	defer func() {
		b.linesChanged.Mark(b.cursor.Line) // after
	}()

	b.cursor.Line -= n
	if b.cursor.Line < 0 {
		b.cursor.Line = 0
	}
	line := b.getCurrentLine()
	if b.cursor.Column > len(line) {
		b.cursor.Column = len(line)
	}
}

// MoveWordLeft moves the cursor left to the previous word
func (b *buffer) MoveWordLeft(locked ...bool) {
	if len(locked) == 0 {
		b.mutex.Lock()
		defer b.mutex.Unlock()
	}

	// if cursor is on the first column, move it to the previous line
	if b.cursor.Column == 0 {
		if b.cursor.Line == 0 {
			return
		}
		b.linesChanged.Mark(b.cursor.Line) // before
		b.cursor.Line--
		b.linesChanged.Mark(b.cursor.Line) // after
		b.cursor.Column = len(b.getCurrentLine())
	}

	// move column by column until previous word is found
	foundWord := false
	line := b.getCurrentLine()
	for colIdx := b.cursor.Column - 1; colIdx >= 0; colIdx-- {
		b.cursor.Column = colIdx
		isPoW := isPartOfWord(line[colIdx])
		if foundWord && (!isPoW || colIdx == 0) {
			if !isPoW {
				b.cursor.Column++
			}
			return
		}
		if isPoW {
			foundWord = true
		}
	}
}

// MoveWordRight moves the cursor right to the next word
func (b *buffer) MoveWordRight(locked ...bool) {
	if len(locked) == 0 {
		b.mutex.Lock()
		defer b.mutex.Unlock()
	}

	// if cursor is on the last column, move to the next line
	foundBreak := false
	idxStartingLine := b.cursor.Line
	if b.cursor.Column == len(b.getCurrentLine()) {
		// if already on the last line, there is nothing to do
		if b.cursor.Line == len(b.lines)-1 {
			return
		}
		b.linesChanged.Mark(b.cursor.Line) // before
		b.cursor.Line++
		b.linesChanged.Mark(b.cursor.Line) // after
		b.cursor.Column = 0
		foundBreak = true
	}

	// go line by line until next word is found
	for lineIdx := b.cursor.Line; lineIdx < len(b.lines); lineIdx++ {
		b.cursor.Line = lineIdx
		b.linesChanged.Mark(b.cursor.Line)
		if lineIdx != idxStartingLine {
			b.cursor.Column = 0
		}

		line := b.lines[lineIdx]
		for colIdx := b.cursor.Column; colIdx < len(line); colIdx++ {
			b.cursor.Column = colIdx
			isPoW := isPartOfWord(line[b.cursor.Column])
			if isPoW && foundBreak {
				return
			}
			if !isPoW {
				foundBreak = true
			}
		}
		b.cursor.Column = len(line)
		foundBreak = true
	}
}

func (b *buffer) NumLines() int {
	return len(b.lines)
}

// Reset resets the buffer to its initial state
func (b *buffer) Reset() {
	b.Set("")
}

// Set overwrites the contents of the buffer with the given string.
func (b *buffer) Set(str string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	str = strings.ReplaceAll(str, "\t", b.tab)

	b.done = false
	b.lines = strings.Split(str, "\n")
	b.linesChanged.MarkAll()
	b.linesRendered = time.Now().Format(time.RFC3339Nano)
	b.cursor = CursorLocation{
		Line:   len(b.lines) - 1,
		Column: len(b.lines[len(b.lines)-1]),
	}
}

// SetTab sets the string to use in place of tab characters.
func (b *buffer) SetTab(tab string) {
	b.tab = tab
}

// String returns the current input from the user.
func (b *buffer) String() string {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return strings.Join(b.lines, "\n")
}

func (b *buffer) getCurrentLine() string {
	return b.getLine(b.cursor.Line)
}

func (b *buffer) getCurrentWord(line string) (string, int, int) {
	if len(line) == 0 || b.cursor.Column >= len(line) || !isPartOfWord(line[b.cursor.Column]) {
		return "", -1, -1
	}

	idxWordStart, idxWordEnd := -1, -1
	for idx := b.cursor.Column; idx >= 0; idx-- {
		if !isPartOfWord(line[idx]) {
			break
		}
		idxWordStart = idx
	}
	for idx := b.cursor.Column; idx < len(line); idx++ {
		if !isPartOfWord(line[idx]) {
			break
		}
		idxWordEnd = idx
	}

	return line[idxWordStart : idxWordEnd+1], idxWordStart, idxWordEnd + 1
}

func (b *buffer) getLine(n int) string {
	return b.lines[n]
}

func (b *buffer) getWordAtCursor(wordDelimiters map[byte]bool) (string, int) {
	line := b.getCurrentLine()
	if b.cursor.Column == len(line) || (b.cursor.Column < len(line) && line[b.cursor.Column] == ' ') {
		idxWordStart := -1
		for idx := b.cursor.Column - 1; idx >= 0; idx-- {
			r := line[idx]
			if wordDelimiters != nil {
				if wordDelimiters[r] {
					break
				}
			} else if !isPartOfWord(r) {
				break
			}
			idxWordStart = idx
		}
		if idxWordStart >= 0 {
			return line[idxWordStart:b.cursor.Column], idxWordStart
		}
	}
	return "", -1
}

type linesChangedMap map[int]bool

func (lc linesChangedMap) Clear() {
	for k := range lc {
		delete(lc, k)
	}
}

func (lc linesChangedMap) AllChanged() bool {
	return lc[-1]
}

func (lc linesChangedMap) AnythingChanged() bool {
	return len(lc) > 0
}

func (lc linesChangedMap) IsChanged(line int) bool {
	return lc[-1] || lc[line]
}

func (lc linesChangedMap) MarkAll() {
	lc[-1] = true
}

func (lc linesChangedMap) Mark(line int) {
	lc[line] = true
}

func (lc linesChangedMap) NothingChanged() bool {
	return len(lc) == 0
}

func (lc linesChangedMap) String() string {
	var lines []int
	for k := range lc {
		lines = append(lines, k)
	}
	sort.Ints(lines)
	return fmt.Sprintf("%v", lines)
}

var (
	nonWordRunes = map[byte]bool{
		' ':  true,
		'(':  true,
		')':  true,
		',':  true,
		'.':  true,
		';':  true,
		'[':  true,
		'\n': true,
		'\t': true,
		']':  true,
		'{':  true,
		'}':  true,
	}
)

func isPartOfWord(r byte) bool {
	return !nonWordRunes[r]
}
