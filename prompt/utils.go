package prompt

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
)

// calculateViewportRange returns the 0-indexed start and stop values given the
// number of total lines, the 0-indexed cursor line and the maximum number of
// lines allowed in the viewport.
func calculateViewportRange(numLines int, cursorLine int, maxHeight int) (int, int) {
	// if max height is not defined, return the entire range
	if maxHeight <= 0 {
		return 0, numLines - 1
	}

	// place the cursor in the middle of the range
	start, stop := cursorLine-(maxHeight/2)-1, cursorLine+(maxHeight/2)
	// expand the window right if start is before range
	for start < 0 || stop-start >= maxHeight {
		start++
	}
	// expand the window left if stop is before end
	for stop-start < maxHeight-1 && stop < numLines-1 {
		stop++
	}
	// move the window left if stop is beyond range
	for stop >= numLines {
		stop--
		if start > 0 {
			start--
		}
	}
	return start, stop
}

// clampValue returns the value that fits inside the min/max range.
func clampValue(val, min, max int) int {
	if min > 0 && val < min {
		val = min
	}
	if max > 0 && val > max {
		val = max
	}
	return val
}

// clampValueAllowZero returns the value that fits inside the min/max range. But
// it allows zero values and returns them without clamping to given range.
func clampValueAllowZero(val, min, max int) int {
	if val == 0 {
		return 0
	}
	return clampValue(val, min, max)
}

func insertCursor(input string, insertIdx int, color Color) string {
	inputRunes := []rune(input)
	colIdx, escSeq, inEscSeq := 0, make([]rune, 0), false
	for idx, r := range inputRunes {
		// skip all the color coding escape sequences
		if r == escSeqStart {
			inEscSeq = true
			escSeq = append(escSeq, r)
			continue
		} else if inEscSeq {
			escSeq = append(escSeq, r)
			if r == escSeqStop {
				inEscSeq = false
			}
			continue
		}
		if string(escSeq) == escSeqReset {
			escSeq = make([]rune, 0)
		}

		if colIdx == insertIdx {
			output := append([]rune{}, inputRunes[:idx]...)
			if len(escSeq) > 0 {
				output = append(output, []rune(escSeqReset)...)
			}
			output = append(output, []rune(color.Sprintf("%c", inputRunes[idx]))...)
			if len(escSeq) > 0 {
				output = append(output, escSeq...)
			}
			output = append(output, inputRunes[idx+1:]...)
			return string(output)
		}
		colIdx++
	}

	return fmt.Sprintf("%s%s", input, color.Sprint(" "))
}

func overwriteContents(input string, newContent string, insertIdx int, maxWidth int) string {
	// if input line is smaller than display width, pad it until it reaches EOL
	inputWidth := text.RuneWidthWithoutEscSequences(input)
	if inputWidth < maxWidth {
		input = text.Pad(input, maxWidth, ' ')
	}

	newContentWidth := text.RuneWidthWithoutEscSequences(newContent)
	// if the new content is longer than allowed width, just trim and return it
	if newContentWidth >= maxWidth {
		return text.Trim(newContent, maxWidth)
	}
	// move the autocomplete dropdown left if it go beyond EOL
	for insertIdx > 0 && insertIdx+newContentWidth > maxWidth {
		insertIdx--
	}

	before := stringSubset(input, 0, insertIdx-1)
	after := stringSubset(input, insertIdx+newContentWidth, inputWidth-1)
	output := fmt.Sprintf("%s%s%s", before, newContent, after)
	return output
}

//gocyclo:ignore
func stringSubset(input string, start, stop int) string {
	if start > stop {
		return ""
	}

	output := make([]rune, 0)
	escSeq, escSeqOpen, inEscSeq := make([]rune, 0), make([]rune, 0), false
	nonEscSeqIdx := -1
	for _, r := range input {
		if r == escSeqStart {
			inEscSeq = true
			escSeq = []rune{r}
		} else if inEscSeq {
			escSeq = append(escSeq, r)
		}
		if !inEscSeq {
			nonEscSeqIdx++
		}

		if !inEscSeq {
			if nonEscSeqIdx >= start && nonEscSeqIdx <= stop {
				if len(escSeq) > 0 {
					output = append(output, escSeq...)

					escSeqOpen = make([]rune, 0)
					if string(escSeq) != escSeqReset {
						escSeqOpen = append([]rune{}, escSeq...)
					}
					escSeq = make([]rune, 0)
				}
				output = append(output, r)
			}
			if nonEscSeqIdx == stop {
				break
			}
		}

		if inEscSeq && r == escSeqStop {
			inEscSeq = false
			if string(escSeq) == escSeqReset {
				if nonEscSeqIdx >= start && nonEscSeqIdx <= stop {
					output = append(output, escSeq...)
				}
				escSeq = make([]rune, 0)
				escSeqOpen = make([]rune, 0)
			}
		}
	}
	if len(escSeqOpen) > 0 {
		output = append(output, []rune(escSeqReset)...)
	}

	return string(output)
}
