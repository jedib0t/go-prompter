package prompt

import (
	"context"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

var (
	testSuggestions = []Suggestion{
		{Value: "auto-complete-1"},
		{Value: "auto-complete-2"},
	}
)

func generateTestPromptWithBuffer(t *testing.T, ctx context.Context, text string, cursor CursorLocation) *prompt {
	p := generateTestPrompt(t, ctx)
	p.buffer.Set(text)
	p.buffer.cursor = cursor
	return p
}

func TestPrompt_handleHistoryExec(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	p := generateTestPrompt(t, ctx)

	output := strings.Builder{}
	p.handleHistoryExec(termenv.NewOutput(&output), 1)
	assert.Contains(t, output.String(), "ERROR: invalid command number: 1")
	assert.Equal(t, "", p.buffer.String())
	assert.False(t, p.buffer.IsDone())
	assert.False(t, p.renderingPaused)

	p.SetHistory(testHistoryCommands)
	output.Reset()
	p.handleHistoryExec(termenv.NewOutput(&output), 1)
	assert.Equal(t, testHistoryCommands[0].Command, p.buffer.String())
	assert.True(t, p.buffer.IsDone())
	assert.False(t, p.renderingPaused)
}

func TestPrompt_handleHistoryList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	p := generateTestPrompt(t, ctx)

	output := strings.Builder{}
	p.handleHistoryList(termenv.NewOutput(&output), 1)
	assert.Contains(t, output.String(), p.history.Render(1, 0))
	assert.Equal(t, "", p.buffer.String())
	assert.False(t, p.buffer.IsDone())
	assert.False(t, p.renderingPaused)
}

func TestPrompt_handleKey(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	t.Run("auto-complete", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.isInAutoComplete = true
		p.suggestions = append(p.suggestions, testSuggestions...)
		p.suggestionsIdx = 0

		output := strings.Builder{}
		err := p.handleKey(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyDown})
		assert.Nil(t, err)
		assert.Equal(t, 1, p.suggestionsIdx)
		assert.Equal(t, "", output.String())
	})

	t.Run("insert", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)

		output := strings.Builder{}
		err := p.handleKey(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
		assert.Nil(t, err)
		assert.Equal(t, "a", p.buffer.String())
	})
}

func TestPrompt_handleKeyAutoComplete(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	t.Run("AutoCompleteChooseNext", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.isInAutoComplete = true
		p.keyMapReversed.AutoComplete[Enter] = AutoCompleteChooseNext
		p.suggestions = append(p.suggestions, testSuggestions...)
		p.suggestionsIdx = 0

		output := strings.Builder{}
		err := p.handleKeyAutoComplete(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, 1, p.suggestionsIdx)
		assert.Equal(t, "", output.String())
	})

	t.Run("AutoCompleteChoosePrevious", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.isInAutoComplete = true
		p.keyMapReversed.AutoComplete[Enter] = AutoCompleteChoosePrevious
		p.suggestions = append(p.suggestions, testSuggestions...)
		p.suggestionsIdx = 1

		output := strings.Builder{}
		err := p.handleKeyAutoComplete(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, 0, p.suggestionsIdx)
		assert.Equal(t, "", output.String())
	})

	t.Run("AutoCompleteSelect", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.buffer.Set("Auto-")
		p.isInAutoComplete = true
		p.keyMapReversed.AutoComplete[Enter] = AutoCompleteSelect
		p.suggestions = append(p.suggestions, testSuggestions...)
		p.suggestionsIdx = 1

		output := strings.Builder{}
		err := p.handleKeyAutoComplete(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "Auto-complete-2 ", p.buffer.String())
		assert.Equal(t, 0, p.suggestionsIdx)
		assert.Equal(t, "", output.String())
	})

	t.Run("fall-through", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.isInAutoComplete = true

		output := strings.Builder{}
		err := p.handleKeyAutoComplete(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
		assert.Nil(t, err)
		assert.Equal(t, "a", p.buffer.String())
		assert.Equal(t, "", output.String())
	})

	t.Run("unknown action", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.isInAutoComplete = true
		p.keyMapReversed.AutoComplete[Enter] = Action("foo")

		output := strings.Builder{}
		err := p.handleKeyAutoComplete(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
	})
}

func TestPrompt_handleKeyInsert(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	testText1 := "test this thing"
	testText2 := "test this thing\nnot this thing"

	t.Run("shortcuts", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.shortcuts = map[KeySequence]string{
			F1: "/help",
		}

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyF1})
		assert.Nil(t, err)
		assert.Equal(t, "/help", p.buffer.String())
		assert.True(t, p.buffer.IsDone())
	})

	t.Run("Abort", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = Abort

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.NotNil(t, err)
		assert.Equal(t, ErrAborted, err)
		assert.Empty(t, p.buffer.String())
	})

	t.Run("AutoComplete", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = AutoComplete

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.True(t, p.forcedAutoComplete())
	})

	t.Run("DeleteCharCurrent", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, "test", CursorLocation{0, 3})
		p.keyMapReversed.Insert[Enter] = DeleteCharCurrent

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "tes", p.buffer.String())
	})

	t.Run("DeleteCharPrevious", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, "test", CursorLocation{0, 3})
		p.keyMapReversed.Insert[Enter] = DeleteCharPrevious

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "tet", p.buffer.String())
	})

	t.Run("DeleteWordNext", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText1, CursorLocation{0, 5})
		p.keyMapReversed.Insert[Enter] = DeleteWordNext

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "test thing", p.buffer.String())
	})

	t.Run("DeleteWordPrevious", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText1, CursorLocation{0, 5})
		p.keyMapReversed.Insert[Enter] = DeleteWordPrevious

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "this thing", p.buffer.String())
	})

	t.Run("EraseEverything", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText1, CursorLocation{0, 5})
		p.keyMapReversed.Insert[Enter] = EraseEverything

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "", p.buffer.String())
	})

	t.Run("EraseToBeginningOfLine", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText1, CursorLocation{0, 5})
		p.keyMapReversed.Insert[Enter] = EraseToBeginningOfLine

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "this thing", p.buffer.String())
	})

	t.Run("EraseToEndOfLine", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText1, CursorLocation{0, 5})
		p.keyMapReversed.Insert[Enter] = EraseToEndOfLine

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "test ", p.buffer.String())
	})

	t.Run("HistoryNext", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText1, CursorLocation{0, 5})
		p.SetHistory(testHistoryCommands)
		p.history.Index = 0
		p.keyMapReversed.Insert[Enter] = HistoryNext

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, testHistoryCommands[1].Command, p.buffer.String())
	})

	t.Run("HistoryPrevious", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText1, CursorLocation{0, 5})
		p.SetHistory(testHistoryCommands)
		p.history.Index = 1
		p.keyMapReversed.Insert[Enter] = HistoryPrevious

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, testHistoryCommands[0].Command, p.buffer.String())
	})

	t.Run("MakeWordCapitalCase", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText1, CursorLocation{0, 5})
		p.keyMapReversed.Insert[Enter] = MakeWordCapitalCase

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "test This thing", p.buffer.String())
	})

	t.Run("MakeWordLowerCase", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText1, CursorLocation{0, 5})
		p.keyMapReversed.Insert[Enter] = MakeWordLowerCase

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "test this thing", p.buffer.String())
	})

	t.Run("MakeWordUpperCase", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText1, CursorLocation{0, 5})
		p.keyMapReversed.Insert[Enter] = MakeWordUpperCase

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "test THIS thing", p.buffer.String())
	})

	t.Run("MoveDownOneLine", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText2, CursorLocation{0, 5})
		p.keyMapReversed.Insert[Enter] = MoveDownOneLine

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 5}, p.buffer.cursor)
	})

	t.Run("MoveLeftOneCharacter", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText2, CursorLocation{1, 5})
		p.keyMapReversed.Insert[Enter] = MoveLeftOneCharacter

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 4}, p.buffer.cursor)
	})

	t.Run("MoveRightOneCharacter", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText2, CursorLocation{1, 5})
		p.keyMapReversed.Insert[Enter] = MoveRightOneCharacter

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 6}, p.buffer.cursor)
	})

	t.Run("MoveUpOneLine", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText2, CursorLocation{1, 5})
		p.keyMapReversed.Insert[Enter] = MoveUpOneLine

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 0, Column: 5}, p.buffer.cursor)
	})

	t.Run("MoveToBeginning", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText2, CursorLocation{1, 14})
		p.keyMapReversed.Insert[Enter] = MoveToBeginning

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 0, Column: 0}, p.buffer.cursor)
	})

	t.Run("MoveToBeginningOfLine", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText2, CursorLocation{1, 14})
		p.keyMapReversed.Insert[Enter] = MoveToBeginningOfLine

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 0}, p.buffer.cursor)
	})

	t.Run("MoveToEnd", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText2, CursorLocation{0, 0})
		p.keyMapReversed.Insert[Enter] = MoveToEnd

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 14}, p.buffer.cursor)
	})

	t.Run("MoveToEndOfLine", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText2, CursorLocation{1, 0})
		p.keyMapReversed.Insert[Enter] = MoveToEndOfLine

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 14}, p.buffer.cursor)
	})

	t.Run("MoveToWordNext", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText2, CursorLocation{1, 0})
		p.keyMapReversed.Insert[Enter] = MoveToWordNext

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 4}, p.buffer.cursor)
	})

	t.Run("MoveToWordPrevious", func(t *testing.T) {
		p := generateTestPromptWithBuffer(t, ctx, testText2, CursorLocation{1, 14})
		p.keyMapReversed.Insert[Enter] = MoveToWordPrevious

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 9}, p.buffer.cursor)
	})

	t.Run("Terminate History Exec", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.SetHistory(testHistoryCommands)
		p.keyMapReversed.Insert[Enter] = Terminate
		p.buffer.InsertString("!1")

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, testHistoryCommands[0].Command, p.buffer.String())
		assert.True(t, p.buffer.IsDone())
		assert.Contains(t, p.debugDataAsString(), "reason=hist.exec")
	})

	t.Run("Terminate History List", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.SetHistory(testHistoryCommands)
		p.keyMapReversed.Insert[Enter] = Terminate
		p.buffer.InsertString("!!")

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "", p.buffer.String())
		assert.False(t, p.buffer.IsDone())
		assert.Contains(t, p.debugDataAsString(), "reason=hist.list")
	})

	t.Run("Terminate Done", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = Terminate
		p.buffer.InsertString(testText1)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, testText1, p.buffer.String())
		assert.True(t, p.buffer.IsDone())
	})

	t.Run("Terminate Insert Newline", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = Terminate
		p.buffer.InsertString(testText1)
		p.SetTerminationChecker(TerminationCheckerSQL())

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "test this thing\n", p.buffer.String())
		assert.False(t, p.buffer.IsDone())
	})

	t.Run("Runes", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.SetTerminationChecker(TerminationCheckerSQL())

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("abc")})
		assert.Nil(t, err)
		assert.Equal(t, "abc", p.buffer.String())
	})

	t.Run("Runes Space", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.SetTerminationChecker(TerminationCheckerSQL())

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeySpace})
		assert.Nil(t, err)
		assert.Equal(t, " ", p.buffer.String())
	})

	t.Run("Runes Space", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.SetTerminationChecker(TerminationCheckerSQL())

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyTab})
		assert.Nil(t, err)
		assert.Equal(t, p.style.TabString, p.buffer.String())
	})

	t.Run("unknown action", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = Action("foo")

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
	})
}
