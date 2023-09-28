package prompt

import (
	"context"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	mock_input "github.com/jedib0t/go-prompter/mocks/input"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPrompt_CursorLocation(t *testing.T) {
	p := &prompt{}
	assert.Equal(t, CursorLocation{}, p.CursorLocation())

	p.buffer = newBuffer()
	assert.Equal(t, CursorLocation{Line: 0, Column: 0}, p.CursorLocation())
	p.buffer.Insert('1')
	assert.Equal(t, CursorLocation{Line: 0, Column: 1}, p.CursorLocation())
	p.buffer.Insert('\n')
	assert.Equal(t, CursorLocation{Line: 1, Column: 0}, p.CursorLocation())
}

func TestPrompt_History(t *testing.T) {
	p := prompt{}
	assert.Empty(t, p.History())

	p.SetHistory(testHistoryCommands)
	assert.Len(t, p.History(), 2)
	for idx, cmd := range p.History() {
		assert.Equal(t, testHistoryCommands[idx].Command, cmd.Command)
		assert.Equal(t, testHistoryCommands[idx].Timestamp, cmd.Timestamp)
	}
}

func TestPrompt_KeyMap(t *testing.T) {
	p := prompt{}
	assert.NotEqual(t, KeyMapDefault, p.keyMap)

	err := p.SetKeyMap(KeyMapDefault)
	assert.Nil(t, err)
	assert.Equal(t, KeyMapDefault.AutoComplete, p.keyMap.AutoComplete)
	assert.Equal(t, KeyMapDefault.Insert, p.keyMap.Insert)
	assert.NotNil(t, p.keyMapReversed)
	if p.keyMapReversed != nil {
		assert.Len(t, p.keyMapReversed.AutoComplete, 3)
		assert.Len(t, p.keyMapReversed.Insert, 28)
	}
}

func TestPrompt_NumLines(t *testing.T) {
	p := prompt{}
	assert.Zero(t, p.NumLines())

	p.buffer = newBuffer()
	assert.Equal(t, 1, p.NumLines())
	p.buffer.Insert('1')
	assert.Equal(t, 1, p.NumLines())
	p.buffer.Insert('\n')
	assert.Equal(t, 2, p.NumLines())
}

func TestPrompt_Prompt(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	chErrors := make(chan error, 1)
	chKeyEvents := make(chan tea.KeyMsg, 1)
	chWindowSizeEvents := make(chan tea.WindowSizeMsg, 1)

	mc := gomock.NewController(t)
	defer mc.Finish()
	mockReader := mock_input.NewMockReader(mc)
	mockReader.EXPECT().Begin(gomock.Any())
	mockReader.EXPECT().Reset()
	mockReader.EXPECT().Errors().AnyTimes().Return(chErrors)
	mockReader.EXPECT().KeyEvents().AnyTimes().Return(chKeyEvents)
	mockReader.EXPECT().WindowSizeEvents().AnyTimes().Return(chWindowSizeEvents)
	mockReader.EXPECT().End()

	p := generateTestPrompt(t, ctx)
	p.reader = mockReader
	go func() {
		time.Sleep(time.Second / 10) // some time for all goroutines to start
		chKeyEvents <- tea.KeyMsg{Type: tea.KeyEscape}
	}()
	userInput, err := p.Prompt(ctx)
	assert.Empty(t, userInput)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAborted, err)
}

func TestPrompt_SetAutoCompleter(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.autoCompleter)

	p.SetAutoCompleter(AutoCompleteGoLangKeywords())
	assert.NotNil(t, p.autoCompleter)
}

func TestPrompt_SetAutoCompleterContextual(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.autoCompleterContextual)

	p.SetAutoCompleterContextual(AutoCompleteGoLangKeywords())
	assert.NotNil(t, p.autoCompleterContextual)
}

func TestPrompt_SetCommandShortcuts(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.shortcuts)

	p.SetCommandShortcuts(map[KeySequence]string{
		Escape: "/quit",
		F1:     "/help",
	})
	assert.Len(t, p.shortcuts, 2)
	assert.Contains(t, p.shortcuts, Escape)
	assert.Contains(t, p.shortcuts, F1)
}

func TestPrompt_SetDebug(t *testing.T) {
	p := prompt{}
	assert.False(t, p.debug)

	p.SetDebug(true)
	assert.True(t, p.debug)

	p.SetDebug(false)
	assert.False(t, p.debug)
}

func TestPrompt_SetHeader(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.headerGenerator)

	p.SetHeader("<title>")
	assert.NotNil(t, p.headerGenerator)
	assert.Equal(t, "<title>", p.headerGenerator(100))
}

func TestPrompt_SetHeaderGenerator(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.headerGenerator)

	p.SetHeaderGenerator(LineRuler(StyleLineNumbersEnabled.Color))
	assert.NotNil(t, p.headerGenerator)
	assert.Equal(t,
		"\x1b[38;5;237;48;5;233m----+----1----+----2----+----3----+----4----+----5----+----6----+----7----+----8\x1b[0m",
		p.headerGenerator(80))
}

func TestPrompt_SetHistory(t *testing.T) {
	p := prompt{}
	assert.Len(t, p.history.Commands, 0)
	assert.Equal(t, 0, p.history.Index)

	p.SetHistory(testHistoryCommands)
	assert.Len(t, p.history.Commands, 2)
	assert.Equal(t, len(testHistoryCommands), p.history.Index)
}

func TestPrompt_SetHistoryExecPrefix(t *testing.T) {
	p := prompt{}
	assert.Empty(t, p.historyExecPrefix)

	p.SetHistoryExecPrefix("!")
	assert.Equal(t, "!", p.historyExecPrefix)
}

func TestPrompt_SetHistoryListPrefix(t *testing.T) {
	p := prompt{}
	assert.Empty(t, p.historyListPrefix)

	p.SetHistoryListPrefix("!!")
	assert.Equal(t, "!!", p.historyListPrefix)
}

func TestPrompt_SetKeyMap(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.keyMapReversed)

	err := p.SetKeyMap(KeyMapDefault)
	assert.Nil(t, err)
	assert.NotNil(t, p.keyMapReversed)
}

func TestPrompt_SetPrefix(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.prefixer)

	p.SetPrefix("> ")
	assert.NotNil(t, p.prefixer)
	assert.Equal(t, "> ", p.prefixer())
}

func TestPrompt_SetPrefixer(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.prefixer)

	p.SetPrefixer(PrefixText("> "))
	assert.NotNil(t, p.prefixer)
	assert.Equal(t, "> ", p.prefixer())
}

func TestPrompt_SetRefreshInterval(t *testing.T) {
	p := prompt{}
	assert.Equal(t, time.Duration(0), p.refreshInterval)

	p.SetRefreshInterval(time.Second)
	assert.Equal(t, time.Second, p.refreshInterval)
}

func TestPrompt_SetStyle(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.style)

	p.SetStyle(StyleDefault)
	assert.NotNil(t, p.style)
	if p.style != nil {
		assert.Equal(t, StyleDefault, *p.style)
	}
}

func TestPrompt_SetSyntaxHighlighter(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.syntaxHighlighter)

	shSQL, err := SyntaxHighlighterSQL()
	assert.NotNil(t, shSQL)
	assert.Nil(t, err)
	if shSQL != nil {
		p.SetSyntaxHighlighter(shSQL)
		assert.NotNil(t, p.syntaxHighlighter)
		assert.Equal(t,
			"\x1b[38;5;81mselect\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;197m*\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;81mfrom\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;231musers\x1b[0m",
			p.syntaxHighlighter("select * from users"))
	}
}

func TestPrompt_SetTerminationChecker(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.terminationChecker)

	p.SetTerminationChecker(TerminationCheckerNone())
	assert.NotNil(t, p.terminationChecker)
	assert.True(t, p.terminationChecker("foo"))

	p.SetTerminationChecker(TerminationCheckerSQL())
	assert.NotNil(t, p.terminationChecker)
	assert.False(t, p.terminationChecker("foo"))
	assert.True(t, p.terminationChecker("/foo"))
	assert.True(t, p.terminationChecker("foo;"))
}

func TestPrompt_SetWidthEnforcer(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.widthEnforcer)

	p.SetWidthEnforcer(WidthEnforcerDefault)
	assert.NotNil(t, p.widthEnforcer)
	assert.Equal(t, "foo\nbar", p.widthEnforcer("foobar", 3))
}

func TestPrompt_Style(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.Style())

	p.SetStyle(StyleDefault)
	assert.NotNil(t, p.Style())
}
