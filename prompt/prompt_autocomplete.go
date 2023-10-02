package prompt

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
)

func (p *prompt) autoComplete(lines []string, cursorPos CursorLocation, startIdx int) []string {
	suggestions, suggestionsIdx := p.getSuggestionsAndIdx()
	if len(suggestions) == 0 {
		p.isInAutoComplete = false
		return lines
	}
	p.isInAutoComplete = true

	// get the line styling
	linePrefix, prefixWidth, _, numLen, _, _ := p.calculateLineStyling(lines)
	word, _ := p.buffer.getWordAtCursor(p.style.AutoComplete.WordDelimiters)
	wordLen := len(word)

	// get the suggestions printed to super-impose on the displayed lines
	suggestionsDropDown := printSuggestionsDropDown(suggestions, suggestionsIdx, p.style.AutoComplete)

	// if the suggestions are going beyond the last line, pad the lines
	numEmptyLinesToAppend := (len(suggestionsDropDown) + 1 + cursorPos.Line - startIdx) - len(lines)
	if numEmptyLinesToAppend > 0 {
		prefix := linePrefix
		if p.style.LineNumbers.Enabled {
			prefix += strings.Repeat("", numLen)
		}
		if prefix != "" {
			prefix += " "
		}
		for numEmptyLinesToAppend > 0 {
			lines = append(lines, prefix)
			numEmptyLinesToAppend--
		}
	}

	displayWidth := p.getDisplayWidth()
	for idx, suggestion := range suggestionsDropDown {
		lineIdx := idx + cursorPos.Line + 1 - startIdx
		lines[lineIdx] = overwriteContents(
			lines[lineIdx], suggestion, prefixWidth+cursorPos.Column-wordLen-1, displayWidth,
		)
	}
	return lines
}

func (p *prompt) updateSuggestions(ctx context.Context) {
	lastLine, lastWord, lastIdx := "", "", -1
	tick := time.Tick(p.refreshInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			if p.isRenderPaused() {
				continue
			}
			lastLine, lastWord, lastIdx = p.updateSuggestionsInternal(lastLine, lastWord, lastIdx)
		}
	}
}

func (p *prompt) updateSuggestionsInternal(lastLine string, lastWord string, lastIdx int) (string, string, int) {
	// grab the current line, word and index
	p.buffer.mutex.Lock()
	line := p.buffer.getCurrentLine()
	location := uint(p.buffer.cursor.Column)
	word, idx := p.buffer.getWordAtCursor(p.style.AutoComplete.WordDelimiters)
	p.buffer.mutex.Unlock()
	minChars := p.style.AutoComplete.MinChars

	// if there is no word currently, clear drop-down
	forced := false
	if p.forcedAutoComplete() {
		forced = true
	} else if word == "" || idx < 0 || (minChars > 0 && len(word) < minChars) {
		p.setSuggestions(make([]Suggestion, 0))
		p.clearDebugData("ac.")
		return line, word, idx
	}

	// if there is no change compared to before, return old result
	p.setDebugData("ac.idx", fmt.Sprint(idx))
	p.setDebugData("ac.word", fmt.Sprintf("%#v", word))
	if (line == lastLine && word == lastWord && idx == lastIdx) && !forced {
		return line, word, idx
	}

	// prep
	var suggestions []Suggestion
	if p.autoCompleterContextual != nil {
		suggestions = append(suggestions, p.autoCompleterContextual(line, word, location)...)
	}
	if p.autoCompleter != nil {
		suggestions = append(suggestions, p.autoCompleter(line, word, location)...)
	}

	// update
	currentSuggestions, _ := p.getSuggestionsAndIdx()
	if fmt.Sprintf("%#v", suggestions) != fmt.Sprintf("%#v", currentSuggestions) {
		p.setSuggestions(suggestions)
	}
	return line, word, idx
}

func printSuggestion(value string, color Color, maxLen int) string {
	if maxLen == 0 {
		return ""
	}

	if len(value) > maxLen {
		value = text.Trim(value, maxLen-1) + "~"
	}
	if len(value) < maxLen {
		value = text.Pad(value, maxLen, ' ')
	}
	value = color.Sprintf(" %s ", value)
	return value
}

func printSuggestionsDropDown(suggestions []Suggestion, suggestionsIdx int, style StyleAutoComplete) []string {
	// calculate the lengths for the values and hints
	lenValue, lenHint := 0, 0
	for _, s := range suggestions {
		if len(s.Value) > lenValue {
			lenValue = len(s.Value)
		}
		if len(s.Hint) > lenHint {
			lenHint = len(s.Hint)
		}
	}
	lenValue = clampValue(lenValue, style.ValueLengthMin, style.ValueLengthMax)
	lenHint = clampValueAllowZero(lenHint, style.HintLengthMin, style.HintLengthMax)
	if suggestionsIdx < 0 {
		suggestionsIdx = 0
	} else if suggestionsIdx >= len(suggestions) {
		suggestionsIdx = len(suggestions) - 1
	}

	// calculate the view port range (range of suggestions to display)
	start, stop := calculateViewportRange(len(suggestions), suggestionsIdx, style.NumItems)

	// generate the scrollbar for the drop-down
	scrollbar, _ := style.Scrollbar.Generate(len(suggestions), suggestionsIdx, style.NumItems)

	// generate the drop-down
	var lines []string
	for idx, s := range suggestions {
		// skip if out of viewport
		if idx < start || idx > stop {
			continue
		}

		valueColor, hintColor := style.ValueColor, style.HintColor
		if idx == suggestionsIdx {
			valueColor, hintColor = style.ValueSelectedColor, style.HintSelectedColor
		}
		value := printSuggestion(s.Value, valueColor, lenValue)
		hint := printSuggestion(s.Hint, hintColor, lenHint)
		scroll := scrollbar[idx-start]
		lines = append(lines, value+hint+scroll)
	}
	return lines
}
