package prompt

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func (p *prompt) handleHistoryExec(output *termenv.Output, cmdNum int) {
	p.pauseRender()
	defer p.resumeRender()

	p.updateModel(false)
	p.renderView(output, "hist.exec", true)

	cmd := p.history.Get(cmdNum - 1)
	if cmd == "" {
		errMsg := fmt.Sprintf("ERROR: invalid command number: %v.\n\n", cmdNum)
		_, _ = output.WriteString(p.style.Colors.Error.Sprintf(errMsg))
		p.buffer.Reset()
		return
	}

	p.buffer.Set(cmd)
	p.buffer.MarkAsDone()
}

func (p *prompt) handleHistoryList(output *termenv.Output, numItems int) {
	p.pauseRender()
	defer p.resumeRender()

	p.updateModel(false)
	p.renderView(output, "hist.list", true)

	_, _ = output.WriteString(p.history.Render(numItems, p.getDisplayWidth()))
	_, _ = output.WriteString("\n")
	p.linesRendered = make([]string, 0)
	p.buffer.Reset()
}

func (p *prompt) handleKey(output *termenv.Output, key tea.KeyMsg) error {
	if p.isInAutoComplete {
		return p.handleKeyAutoComplete(output, key)
	} else {
		return p.handleKeyInsert(output, key)
	}
}

func (p *prompt) handleKeyAutoComplete(output *termenv.Output, key tea.KeyMsg) error {
	action := p.translateKeyToAutoCompleteAction(key)
	switch action {
	case AutoCompleteChooseNext:
		if p.changeSuggestionsIdx(1) {
			p.updateModel(true)
		}
	case AutoCompleteChoosePrevious:
		if p.changeSuggestionsIdx(-1) {
			p.updateModel(true)
		}
	case AutoCompleteSelect:
		word, _ := p.buffer.getWordAtCursor()
		suggestions, suggestionsIdx := p.getSuggestionsAndIdx()
		if suggestionsIdx < len(suggestions) {
			suggestion := suggestions[suggestionsIdx].Value
			if !strings.HasPrefix(suggestion, word) { // case mismatch
				word = strings.ToLower(word) // case-insensitive auto-complete converts everything to lower-case
			}
			suggestion = strings.Replace(suggestion, word, "", 1)
			p.buffer.InsertString(suggestion + " ")
			p.forceAutoComplete(false)
			p.setSuggestionsIdx(0)
		}
	default:
		return p.handleKeyInsert(output, key)
	}
	return nil
}

//gocyclo:ignore
func (p *prompt) handleKeyInsert(output *termenv.Output, key tea.KeyMsg) error {
	ks := p.translateKeyToKeySequence(key)
	if shortcut, ok := p.shortcuts[ks]; ok {
		p.buffer.Set(shortcut)
		p.buffer.MarkAsDone()
		return nil
	}

	action := p.translateKeyToInsertAction(key)
	switch action {
	case Abort:
		return ErrAborted
	case AutoComplete:
		p.forceAutoComplete(true)
	case DeleteCharCurrent:
		p.buffer.DeleteForward(1)
	case DeleteCharPrevious:
		p.buffer.DeleteBackward(1)
	case DeleteWordNext:
		p.buffer.DeleteWordForward()
	case DeleteWordPrevious:
		p.buffer.DeleteWordBackward()
	case EraseEverything:
		p.buffer.Reset()
	case EraseToBeginningOfLine:
		p.buffer.DeleteBackwardToBeginningOfLine()
	case EraseToEndOfLine:
		p.buffer.DeleteForwardToEndOfLine()
	case HistoryNext:
		p.buffer.Set(p.history.GetNext())
	case HistoryPrevious:
		p.buffer.Set(p.history.GetPrev())
	case MakeWordCapitalCase:
		p.buffer.MakeWordCapitalCase()
	case MakeWordLowerCase:
		p.buffer.MakeWordLowerCase()
	case MakeWordUpperCase:
		p.buffer.MakeWordUpperCase()
	case MoveDownOneLine:
		p.buffer.MoveDown(1)
	case MoveLeftOneCharacter:
		p.buffer.MoveLeft(1)
	case MoveRightOneCharacter:
		p.buffer.MoveRight(1)
	case MoveUpOneLine:
		p.buffer.MoveUp(1)
	case MoveToBeginning:
		p.buffer.MoveToBeginning()
	case MoveToBeginningOfLine:
		p.buffer.MoveToBeginningOfLine()
	case MoveToEnd:
		p.buffer.MoveToEnd()
	case MoveToEndOfLine:
		p.buffer.MoveToEndOfLine()
	case MoveToWordNext:
		p.buffer.MoveWordRight()
	case MoveToWordPrevious:
		p.buffer.MoveWordLeft()
	case Terminate:
		input := p.buffer.String()
		if histCmd := p.processHistoryCommand(input); histCmd.Type != historyCommandNone {
			switch histCmd.Type {
			case historyCommandExec:
				p.handleHistoryExec(output, histCmd.Value)
			case historyCommandList:
				p.handleHistoryList(output, histCmd.Value)
			}
		} else if p.terminationChecker(input) {
			p.buffer.MarkAsDone()
		} else {
			p.buffer.Insert('\n')
		}
	default:
		if key.Type == tea.KeyRunes {
			for _, r := range key.Runes {
				p.buffer.Insert(r)
			}
		} else if key.Type == tea.KeySpace {
			p.buffer.Insert(' ')
		} else if key.Type == tea.KeyTab {
			p.buffer.InsertString(p.style.TabString)
		}
	}
	return nil
}
