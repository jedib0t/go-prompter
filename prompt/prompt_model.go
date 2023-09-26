package prompt

import (
	"fmt"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
)

func (p *prompt) updateModel(isBeingEdited bool) {
	lines, cursorPos := p.buffer.Display()

	timeStart := time.Now()
	var linesToRender []string

	// header
	if header := p.getHeader(); header != "" {
		for _, line := range strings.Split(header, "\n") {
			linesToRender = append(linesToRender, line)
		}
	}

	// syntax highlight
	timeSyntaxStart := time.Now()
	if p.syntaxHighlighter != nil {
		lines = p.doSyntaxHighlighting(lines)
	}
	timeSyntax := time.Since(timeSyntaxStart)

	// render the input lines
	timeBufferStart := time.Now()
	linesFromBuffer, startIdx := p.generateModelLines(lines, cursorPos, isBeingEdited)
	timeBuffer := time.Since(timeBufferStart)

	// auto-complete
	timeAutoCompleteStart := time.Now()
	if isBeingEdited {
		linesToRender = append(linesToRender, p.autoComplete(linesFromBuffer, cursorPos, startIdx)...)
	} else {
		linesToRender = append(linesToRender, linesFromBuffer...)
	}
	timeAutoComplete := time.Since(timeAutoCompleteStart)

	// footer
	if footer := p.getFooter(); footer != "" {
		for _, line := range strings.Split(footer, "\n") {
			linesToRender = append(linesToRender, line)
		}
	}

	p.linesMutex.Lock()
	p.linesToRender = linesToRender
	p.timeGen = time.Since(timeStart).Round(time.Microsecond)
	p.timeSyntaxGen = timeSyntax.Round(time.Microsecond)
	p.timeBufferGen = timeBuffer.Round(time.Microsecond)
	p.timeAutoComplete = timeAutoComplete.Round(time.Microsecond)
	p.linesMutex.Unlock()
}

func (p *prompt) calculateLineStyling(lines []string) (prefix string, prefixWidth int, numColor Color, numLen int, numFmt string, numNone string) {
	// get the line prefix
	if p.prefixer != nil {
		prefix = p.prefixer()
		if prefix != "" {
			prefixWidth += text.RuneWidthWithoutEscSequences(prefix)
		}
	}

	// if enabled, get the lines number styling info
	if p.style.LineNumbers.Enabled {
		numDigits := len(fmt.Sprint(len(lines)))
		zeroPrefix := ""
		if p.style.LineNumbers.ZeroPrefixed {
			zeroPrefix = "0"
		}

		numColor = p.style.LineNumbers.Color
		numLen = 1 + numDigits + 1 // with padding
		numFmt = fmt.Sprintf(" %%%s%dd ", zeroPrefix, numDigits)
		numNone = numColor.Sprintf(fmt.Sprintf(" %%%ds ", numDigits), " ")

		prefixWidth += text.RuneWidthWithoutEscSequences(fmt.Sprintf(numFmt, 1))
		prefixWidth += 1 // margin
	}

	return
}

//gocyclo:ignore
func (p *prompt) generateModelLines(lines []string, cursorPos CursorLocation, isBeingEdited bool) ([]string, int) {
	// get the line styling
	linePrefix, prefixWidth, lineNumColor, _, lineNumFmt, lineNumNone := p.calculateLineStyling(lines)

	// restrict number of lines rendered if a max-height was set
	start, stop := calculateViewportRange(len(lines), cursorPos.Line, int(p.style.Dimensions.HeightMax))
	//if p.debug {
	//	p.debugData["lines"] = fmt.Sprintf("%d-%d", start, stop)
	//}
	scrollbar, isScrollBarVisible := p.style.Scrollbar.Generate(
		len(lines), cursorPos.Line, int(p.style.Dimensions.HeightMax),
	)

	// calculate remaining width for actual content
	remainingWidth := p.getDisplayWidth() - prefixWidth
	if isScrollBarVisible {
		remainingWidth -= 1
	}

	// render the lines
	linesOut := make([]string, 0)
	for lineIdx, line := range lines {
		// skip if out of viewport
		if lineIdx < start || lineIdx > stop {
			continue
		}

		// insert cursor
		if isBeingEdited && lineIdx == cursorPos.Line && p.style.Cursor.Enabled {
			line = insertCursor(line, cursorPos.Column, p.getCursorColor())
		}

		// split line into multiple lines if longer than viewport width
		subLines := []string{line}
		if text.RuneWidthWithoutEscSequences(line) > remainingWidth {
			subLines = strings.Split(p.widthEnforcer(line, remainingWidth), "\n")
		}
		for subLineIdx, subLine := range subLines {
			out := strings.Builder{}
			if linePrefix != "" {
				_, _ = out.WriteString(linePrefix)
			}
			if p.style.LineNumbers.Enabled {
				if subLineIdx > 0 { // content continues into next physical line
					_, _ = out.WriteString(lineNumNone)
				} else {
					_, _ = out.WriteString(lineNumColor.Sprintf(lineNumFmt, lineIdx+1))
				}
				_, _ = out.WriteString(" ") // margin
			}
			if isScrollBarVisible {
				subLine = text.Pad(subLine, remainingWidth, ' ')
			}
			_, _ = out.WriteString(fmt.Sprintf("%s", subLine))
			if isScrollBarVisible {
				_, _ = out.WriteString(scrollbar[lineIdx-start])
			}

			linesOut = append(linesOut, out.String())
		}
	}

	// add empty lines if number of lines is less than minimum height
	for p.style.Dimensions.HeightMin > 0 && len(linesOut) < int(p.style.Dimensions.HeightMin) {
		if p.style.LineNumbers.Enabled {
			linesOut = append(linesOut, linePrefix+lineNumNone)
		} else {
			linesOut = append(linesOut, "")
		}
	}

	return linesOut, start
}
