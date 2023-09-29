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

type actionHandler func(p *prompt, output *termenv.Output, key tea.KeyMsg) error

var autoCompleteActionHandlerMap = map[Action]actionHandler{
	AutoCompleteChooseNext: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		if p.changeSuggestionsIdx(1) {
			p.updateModel(true)
		}
		return nil
	},
	AutoCompleteChoosePrevious: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		if p.changeSuggestionsIdx(-1) {
			p.updateModel(true)
		}
		return nil
	},
	AutoCompleteSelect: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
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
		return nil
	},
	None: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		return p.handleKeyInsert(output, key)
	},
}

var insertActionHandlerMap = map[Action]actionHandler{
	Abort: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		return ErrAborted
	},
	AutoComplete: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.forceAutoComplete(true)
		return nil
	},
	DeleteCharCurrent: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.DeleteForward(1)
		return nil
	},
	DeleteCharPrevious: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.DeleteBackward(1)
		return nil
	},
	DeleteWordNext: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.DeleteWordForward()
		return nil
	},
	DeleteWordPrevious: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.DeleteWordBackward()
		return nil
	},
	EraseEverything: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.Reset()
		return nil
	},
	EraseToBeginningOfLine: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.DeleteBackwardToBeginningOfLine()
		return nil
	},
	EraseToEndOfLine: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.DeleteForwardToEndOfLine()
		return nil
	},
	HistoryNext: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.Set(p.history.GetNext())
		p.resetSuggestions()
		return nil
	},
	HistoryPrevious: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.Set(p.history.GetPrev())
		p.resetSuggestions()
		return nil
	},
	MakeWordCapitalCase: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.MakeWordCapitalCase()
		return nil
	},
	MakeWordLowerCase: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.MakeWordLowerCase()
		return nil
	},
	MakeWordUpperCase: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.MakeWordUpperCase()
		return nil
	},
	MoveDownOneLine: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.MoveDown(1)
		return nil
	},
	MoveLeftOneCharacter: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.MoveLeft(1)
		return nil
	},
	MoveRightOneCharacter: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.MoveRight(1)
		return nil
	},
	MoveUpOneLine: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.MoveUp(1)
		return nil
	},
	MoveToBeginning: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.MoveToBeginning()
		return nil
	},
	MoveToBeginningOfLine: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.MoveToBeginningOfLine()
		return nil
	},
	MoveToEnd: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.MoveToEnd()
		return nil
	},
	MoveToEndOfLine: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.MoveToEndOfLine()
		return nil
	},
	MoveToWordNext: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.MoveWordRight()
		return nil
	},
	MoveToWordPrevious: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		p.buffer.MoveWordLeft()
		return nil
	},
	Terminate: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
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
		p.forceAutoComplete(false)
		p.resetSuggestions()
		return nil
	},
	None: func(p *prompt, output *termenv.Output, key tea.KeyMsg) error {
		if key.Type == tea.KeyRunes {
			for _, r := range key.Runes {
				p.buffer.Insert(r)
			}
			p.forceAutoComplete(false)
			p.resetSuggestions()
		} else if key.Type == tea.KeySpace {
			p.buffer.Insert(' ')
			p.forceAutoComplete(false)
			p.resetSuggestions()
		} else if key.Type == tea.KeyTab {
			p.buffer.InsertString(p.style.TabString)
			p.forceAutoComplete(false)
			p.resetSuggestions()
		}
		return nil
	},
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
	handler, ok := autoCompleteActionHandlerMap[action]
	if ok && handler != nil {
		p.setDebugData("action", string(action))
		return handler(p, output, key)
	}
	return nil
}

func (p *prompt) handleKeyInsert(output *termenv.Output, key tea.KeyMsg) error {
	ks := translateKeyToKeySequence(key)
	if shortcut, ok := p.shortcuts[ks]; ok {
		p.buffer.Set(shortcut)
		p.buffer.MarkAsDone()
		return nil
	}

	action := p.translateKeyToInsertAction(key)
	handler, ok := insertActionHandlerMap[action]
	if ok && handler != nil {
		p.setDebugData("action", string(action))
		return handler(p, output, key)
	}
	return nil
}
