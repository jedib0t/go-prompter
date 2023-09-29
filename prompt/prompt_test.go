package prompt

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	mock_input "github.com/jedib0t/go-prompter/mocks/input"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	errFoo = errors.New("test-error-foo")
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

func TestPrompt_IsActive(t *testing.T) {
	p := prompt{}
	assert.False(t, p.IsActive())

	p.markActive()
	assert.True(t, p.IsActive())

	p.markInactive()
	assert.False(t, p.IsActive())
}

func TestPrompt_KeyMap(t *testing.T) {
	p := prompt{}
	assert.NotEqual(t, KeyMapDefault, p.keyMap)

	err := p.SetKeyMap(KeyMapDefault)
	assert.Nil(t, err)
	assert.Equal(t, KeyMapDefault.AutoComplete, p.KeyMap().AutoComplete)
	assert.Equal(t, KeyMapDefault.Insert, p.KeyMap().Insert)
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

	t.Run("style error", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)

		s := StyleDefault
		s.Dimensions.WidthMin = 50
		s.Dimensions.WidthMax = 40
		p.SetStyle(s)

		userInput, err := p.Prompt(ctx)
		assert.Empty(t, userInput)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, ErrInvalidDimensions))
	})

	t.Run("input reader error", func(t *testing.T) {
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
			<-time.After(time.Second / 10) // some time for all goroutines to start
			chErrors <- fmt.Errorf("test-error")
		}()
		userInput, err := p.Prompt(ctx)
		assert.Empty(t, userInput)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "test-error")
	})

	t.Run("input error", func(t *testing.T) {
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
			<-time.After(time.Second / 10) // some time for all goroutines to start
			chKeyEvents <- tea.KeyMsg{Type: tea.KeyCtrlC}
		}()
		userInput, err := p.Prompt(ctx)
		assert.Empty(t, userInput)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, ErrAborted))
	})

	t.Run("no error", func(t *testing.T) {
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
			<-time.After(time.Second / 10) // some time for all goroutines to start
			chKeyEvents <- tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("abc\tdef")}
			<-time.After(time.Second / 4) // for some rendering
			chWindowSizeEvents <- tea.WindowSizeMsg{Width: int(p.style.Dimensions.WidthMax + 20)}
			<-time.After(time.Second / 4) // for some rendering
			chKeyEvents <- tea.KeyMsg{Type: tea.KeyEnter}
		}()
		userInput, err := p.Prompt(ctx)
		assert.Equal(t, "abc    def", userInput)
		assert.Nil(t, err)
		assert.Equal(t, int(p.style.Dimensions.WidthMax), p.displayWidth)

		out, ok := p.output.(*strings.Builder)
		if ok {
			outString := out.String()
			t.Log(outString)
			assert.Contains(t, outString, "TestPrompt_Prompt/no_error")
			assert.Contains(t, outString, "abc    def")
		}
	})
}

func TestPrompt_SendInput(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	t.Run("send KeySequence error", func(t *testing.T) {
		mc := gomock.NewController(t)
		defer mc.Finish()
		mockReader := mock_input.NewMockReader(mc)
		mockReader.EXPECT().Send(keySequenceKeyMsgMap[F1]).
			Return(errFoo)

		p := generateTestPrompt(t, ctx)
		p.reader = mockReader
		err := p.SendInput([]any{F1}, time.Second)
		assert.NotNil(t, err)
		assert.Equal(t, errFoo, err)
	})

	t.Run("send KeySequence", func(t *testing.T) {
		mc := gomock.NewController(t)
		defer mc.Finish()
		mockReader := mock_input.NewMockReader(mc)
		mockReader.EXPECT().Send(keySequenceKeyMsgMap[F1])

		p := generateTestPrompt(t, ctx)
		p.reader = mockReader
		err := p.SendInput([]any{F1}, time.Second)
		assert.Nil(t, err)
	})

	t.Run("send time.Duration", func(t *testing.T) {
		mc := gomock.NewController(t)
		defer mc.Finish()

		p := generateTestPrompt(t, ctx)
		err := p.SendInput([]any{time.Microsecond})
		assert.Nil(t, err)
	})

	t.Run("send rune error", func(t *testing.T) {
		mc := gomock.NewController(t)
		defer mc.Finish()
		mockReader := mock_input.NewMockReader(mc)
		mockReader.EXPECT().Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}).
			Return(errFoo)

		p := generateTestPrompt(t, ctx)
		p.reader = mockReader
		err := p.SendInput([]any{'a'})
		assert.NotNil(t, err)
		assert.Equal(t, errFoo, err)
	})

	t.Run("send rune", func(t *testing.T) {
		mc := gomock.NewController(t)
		defer mc.Finish()
		mockReader := mock_input.NewMockReader(mc)
		mockReader.EXPECT().Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

		p := generateTestPrompt(t, ctx)
		p.reader = mockReader
		err := p.SendInput([]any{'a'})
		assert.Nil(t, err)
	})

	t.Run("send string error", func(t *testing.T) {
		mc := gomock.NewController(t)
		defer mc.Finish()
		mockReader := mock_input.NewMockReader(mc)
		mockReader.EXPECT().Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}).
			Return(errFoo)

		p := generateTestPrompt(t, ctx)
		p.reader = mockReader
		err := p.SendInput([]any{"abc"})
		assert.NotNil(t, err)
		assert.Equal(t, errFoo, err)
	})

	t.Run("send string/[]rune", func(t *testing.T) {
		mc := gomock.NewController(t)
		defer mc.Finish()

		for _, input := range []any{"abc", []rune("abc")} {
			mockReader := mock_input.NewMockReader(mc)
			gomock.InOrder(
				mockReader.EXPECT().Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}),
				mockReader.EXPECT().Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}}),
				mockReader.EXPECT().Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}),
			)

			p := generateTestPrompt(t, ctx)
			p.reader = mockReader
			err := p.SendInput([]any{input})
			assert.Nil(t, err)

		}
	})

	t.Run("send unsupported", func(t *testing.T) {
		mc := gomock.NewController(t)
		defer mc.Finish()

		p := generateTestPrompt(t, ctx)
		err := p.SendInput([]any{p})
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, ErrUnsupportedInput))
	})
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

func TestPrompt_SetFooter(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.footerGenerator)

	p.SetFooter("<title>")
	assert.NotNil(t, p.footerGenerator)
	assert.Equal(t, "<title>", p.footerGenerator(100))
}

func TestPrompt_SetFooterGenerator(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.footerGenerator)

	p.SetFooterGenerator(LineRuler(StyleLineNumbersEnabled.Color))
	assert.NotNil(t, p.footerGenerator)
	assert.Equal(t,
		"\x1b[38;5;240;48;5;236m----+----1----+----2----+----3----+----4----+----5----+----6----+----7----+----8\x1b[0m",
		p.footerGenerator(80))
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
		"\x1b[38;5;240;48;5;236m----+----1----+----2----+----3----+----4----+----5----+----6----+----7----+----8\x1b[0m",
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

func TestPrompt_SetInput(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.input)
	assert.Equal(t, os.Stdin, p.getInputReader())

	r := bytes.NewReader([]byte("test"))
	p.SetInput(r)
	assert.NotNil(t, p.input)
	assert.Equal(t, r, p.getInputReader())
	assert.NotNil(t, p.reader)

	p.SetInput(nil)
	assert.Nil(t, p.input)
	assert.Equal(t, os.Stdin, p.getInputReader())
	assert.NotNil(t, p.reader)
}

func TestPrompt_SetKeyMap(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.keyMapReversed)

	err := p.SetKeyMap(KeyMapDefault)
	assert.Nil(t, err)
	assert.NotNil(t, p.keyMapReversed)

	km := KeyMapDefault
	km.Insert.MoveToEnd = km.Insert.MoveToBeginning
	err = p.SetKeyMap(km)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrDuplicateKeyAssignment))
}

func TestPrompt_SetOutput(t *testing.T) {
	p := prompt{}
	assert.Nil(t, p.output)
	assert.Equal(t, os.Stdout, p.getOutputWriter())

	w := &strings.Builder{}
	p.SetOutput(w)
	assert.NotNil(t, p.output)
	assert.Equal(t, w, p.getOutputWriter())

	p.SetOutput(nil)
	assert.Nil(t, p.output)
	assert.Equal(t, os.Stdout, p.getOutputWriter())
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

	p.SetRefreshInterval(0)
	assert.Equal(t, DefaultRefreshInterval, p.refreshInterval)
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

func TestPrompt_changeSuggestionsIdx(t *testing.T) {
	p := prompt{}
	p.suggestions = testSuggestions
	p.suggestionsIdx = 0

	assert.False(t, p.changeSuggestionsIdx(0))
	assert.Equal(t, 0, p.suggestionsIdx)

	assert.True(t, p.changeSuggestionsIdx(1))
	assert.Equal(t, 1, p.suggestionsIdx)

	assert.False(t, p.changeSuggestionsIdx(1))
	assert.Equal(t, 1, p.suggestionsIdx)

	assert.True(t, p.changeSuggestionsIdx(-1))
	assert.Equal(t, 0, p.suggestionsIdx)

	assert.False(t, p.changeSuggestionsIdx(-1))
	assert.Equal(t, 0, p.suggestionsIdx)
}

func TestPrompt_clearDebugData(t *testing.T) {
	p := prompt{}
	p.debugData = make(map[string]string)
	assert.Equal(t, "n/a", p.debugDataAsString())

	p.setDebugData("k1", "v1")
	p.setDebugData("k2", "v2")
	p.setDebugData("v1", "v1")
	p.setDebugData("v2", "v1")
	assert.Len(t, p.debugData, 4)
	p.clearDebugData("k")
	assert.Len(t, p.debugData, 2)
	p.clearDebugData("v")
	assert.Len(t, p.debugData, 0)

	p.setDebugData("k1", "v1")
	p.setDebugData("k2", "v2")
	p.setDebugData("v1", "v1")
	p.setDebugData("v2", "v1")
	assert.Len(t, p.debugData, 4)
	p.clearDebugData()
	assert.Len(t, p.debugData, 0)
}

func TestPrompt_updateCursorColors(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p := prompt{}
	p.SetStyle(StyleDefault)
	p.Style().Cursor.Blink = true
	p.Style().Cursor.BlinkInterval = time.Millisecond * 20
	p.setCursorColor(p.Style().Cursor.Color)

	go p.updateCursorColors(ctx)
	assert.Equal(t, p.Style().Cursor.Color, p.getCursorColor())
	<-time.After(time.Millisecond * 25)
	assert.Equal(t, p.Style().Cursor.ColorAlt, p.getCursorColor())
	<-time.After(time.Millisecond * 25)
	assert.Equal(t, p.Style().Cursor.Color, p.getCursorColor())
}

func TestPrompt_updateDisplayWidth(t *testing.T) {
	p := prompt{}
	assert.Equal(t, 0, p.displayWidth)
	widthMin, widthMax := 80, 120

	p.SetStyle(StyleDefault)
	p.Style().Dimensions.WidthMin = uint(widthMin)
	p.Style().Dimensions.WidthMax = uint(widthMax)
	p.updateDisplayWidth(5)
	assert.Equal(t, widthMin, p.displayWidth)

	p.updateDisplayWidth(80)
	assert.Equal(t, widthMin, p.displayWidth)

	p.updateDisplayWidth(100)
	assert.Equal(t, 100, p.displayWidth)

	p.updateDisplayWidth(120)
	assert.Equal(t, widthMax, p.displayWidth)

	p.updateDisplayWidth(150)
	assert.Equal(t, widthMax, p.displayWidth)

	p.debug = true
	p.updateDisplayWidth(150)
	assert.Equal(t, widthMax-debugMarginWidth, p.displayWidth)
}

func TestPrompt_updateHeaderAndFooterAsync(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p := prompt{}
	p.SetRefreshInterval(time.Millisecond * 20)
	p.SetHeaderGenerator(LineSimple("test-header"))
	p.displayWidth = 120

	go p.updateHeaderAndFooterAsync(ctx)
	assert.Equal(t, "", p.getHeader())
	<-time.After(time.Millisecond * 50)
	assert.Equal(t, "test-header", p.getHeader())
}

func Test_translateKeyToKeySequence(t *testing.T) {
	assert.Equal(t, AltA, translateKeyToKeySequence(tea.KeyMsg{
		Alt: true, Runes: []rune{'a'},
	}))
	assert.Equal(t, AltB, translateKeyToKeySequence(tea.KeyMsg{
		Alt: true, Runes: []rune{'b'},
	}))
	assert.Equal(t, Enter, translateKeyToKeySequence(tea.KeyMsg{
		Type: tea.KeyEnter,
	}))
}
