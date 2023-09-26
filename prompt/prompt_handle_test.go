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
		p.suggestions = []Suggestion{
			{Value: "auto-complete-1"},
			{Value: "auto-complete-2"},
		}
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
		p.suggestions = []Suggestion{
			{Value: "auto-complete-1"},
			{Value: "auto-complete-2"},
		}
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
		p.suggestions = []Suggestion{
			{Value: "auto-complete-1"},
			{Value: "auto-complete-2"},
		}
		p.suggestionsIdx = 1

		output := strings.Builder{}
		err := p.handleKeyAutoComplete(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, 0, p.suggestionsIdx)
		assert.Equal(t, "", output.String())
	})

	t.Run("AutoCompleteSelect", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.isInAutoComplete = true
		p.keyMapReversed.AutoComplete[Enter] = AutoCompleteSelect
		p.suggestions = []Suggestion{
			{Value: "auto-complete-1"},
			{Value: "auto-complete-2"},
		}
		p.suggestionsIdx = 1

		output := strings.Builder{}
		err := p.handleKeyAutoComplete(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "auto-complete-2 ", p.buffer.String())
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
}

func TestPrompt_handleKeyInsert(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

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

	t.Run("DeleteCharCurrent", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = DeleteCharCurrent
		p.buffer.InsertString("test")
		p.buffer.MoveLeft(1)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "tes", p.buffer.String())
	})

	t.Run("DeleteCharPrevious", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = DeleteCharPrevious
		p.buffer.InsertString("test")
		p.buffer.MoveLeft(1)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "tet", p.buffer.String())
	})

	t.Run("DeleteWordNext", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = DeleteWordNext
		p.buffer.InsertString("test this thing")
		p.buffer.MoveToBeginning()
		p.buffer.MoveRight(5)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "test thing", p.buffer.String())
	})

	t.Run("DeleteWordPrevious", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = DeleteWordPrevious
		p.buffer.InsertString("test this thing")
		p.buffer.MoveToBeginning()
		p.buffer.MoveRight(5)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "this thing", p.buffer.String())
	})

	t.Run("EraseEverything", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = EraseEverything
		p.buffer.InsertString("test this thing")
		p.buffer.MoveToBeginning()
		p.buffer.MoveRight(5)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "", p.buffer.String())
	})

	t.Run("EraseToBeginningOfLine", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = EraseToBeginningOfLine
		p.buffer.InsertString("test this thing")
		p.buffer.MoveToBeginning()
		p.buffer.MoveRight(5)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "this thing", p.buffer.String())
	})

	t.Run("EraseToEndOfLine", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = EraseToEndOfLine
		p.buffer.InsertString("test this thing")
		p.buffer.MoveToBeginning()
		p.buffer.MoveRight(5)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "test ", p.buffer.String())
	})

	t.Run("HistoryNext", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.SetHistory(testHistoryCommands)
		p.history.Index = 0
		p.keyMapReversed.Insert[Enter] = HistoryNext
		p.buffer.InsertString("test this thing")
		p.buffer.MoveToBeginning()
		p.buffer.MoveRight(5)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, testHistoryCommands[1].Command, p.buffer.String())
	})

	t.Run("HistoryPrevious", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.SetHistory(testHistoryCommands)
		p.history.Index = 1
		p.keyMapReversed.Insert[Enter] = HistoryPrevious
		p.buffer.InsertString("test this thing")
		p.buffer.MoveToBeginning()
		p.buffer.MoveRight(5)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, testHistoryCommands[0].Command, p.buffer.String())
	})

	t.Run("MakeWordCapitalCase", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = MakeWordCapitalCase
		p.buffer.InsertString("test this thing")
		p.buffer.MoveToBeginning()
		p.buffer.MoveRight(5)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "test This thing", p.buffer.String())
	})

	t.Run("MakeWordLowerCase", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = MakeWordLowerCase
		p.buffer.InsertString("test THIS thing")
		p.buffer.MoveToBeginning()
		p.buffer.MoveRight(5)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "test this thing", p.buffer.String())
	})

	t.Run("MakeWordUpperCase", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = MakeWordUpperCase
		p.buffer.InsertString("test this thing")
		p.buffer.MoveToBeginning()
		p.buffer.MoveRight(5)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "test THIS thing", p.buffer.String())
	})

	t.Run("MoveDownOneLine", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = MoveDownOneLine
		p.buffer.InsertString("test this thing\nnot this thing")
		p.buffer.MoveToBeginning()
		p.buffer.MoveRight(5)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 5}, p.buffer.cursor)
	})

	t.Run("MoveLeftOneCharacter", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = MoveLeftOneCharacter
		p.buffer.InsertString("test this thing\nnot this thing")
		p.buffer.MoveToBeginning()
		p.buffer.MoveRight(5)
		p.buffer.MoveDown(1)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 4}, p.buffer.cursor)
	})

	t.Run("MoveRightOneCharacter", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = MoveRightOneCharacter
		p.buffer.InsertString("test this thing\nnot this thing")
		p.buffer.MoveToBeginning()
		p.buffer.MoveRight(5)
		p.buffer.MoveDown(1)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 6}, p.buffer.cursor)
	})

	t.Run("MoveUpOneLine", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = MoveUpOneLine
		p.buffer.InsertString("test this thing\nnot this thing")
		p.buffer.MoveToBeginning()
		p.buffer.MoveRight(5)
		p.buffer.MoveDown(1)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 0, Column: 5}, p.buffer.cursor)
	})

	t.Run("MoveToBeginning", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = MoveToBeginning
		p.buffer.InsertString("test this thing\nnot this thing")
		assert.Equal(t, CursorLocation{Line: 1, Column: 14}, p.buffer.cursor)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 0, Column: 0}, p.buffer.cursor)
	})

	t.Run("MoveToBeginningOfLine", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = MoveToBeginningOfLine
		p.buffer.InsertString("test this thing\nnot this thing")
		assert.Equal(t, CursorLocation{Line: 1, Column: 14}, p.buffer.cursor)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 0}, p.buffer.cursor)
	})

	t.Run("MoveToEnd", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = MoveToEnd
		p.buffer.InsertString("test this thing\nnot this thing")
		p.buffer.MoveToBeginning()
		assert.Equal(t, CursorLocation{Line: 0, Column: 0}, p.buffer.cursor)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 14}, p.buffer.cursor)
	})

	t.Run("MoveToEndOfLine", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = MoveToEndOfLine
		p.buffer.InsertString("test this thing\nnot this thing")
		p.buffer.MoveToBeginningOfLine()
		assert.Equal(t, CursorLocation{Line: 1, Column: 0}, p.buffer.cursor)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 14}, p.buffer.cursor)
	})

	t.Run("MoveToWordNext", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = MoveToWordNext
		p.buffer.InsertString("test this thing\nnot this thing")
		p.buffer.MoveToBeginningOfLine()
		assert.Equal(t, CursorLocation{Line: 1, Column: 0}, p.buffer.cursor)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 4}, p.buffer.cursor)
	})

	t.Run("MoveToWordPrevious", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = MoveToWordPrevious
		p.buffer.InsertString("test this thing\nnot this thing")
		p.buffer.MoveToEndOfLine()
		assert.Equal(t, CursorLocation{Line: 1, Column: 14}, p.buffer.cursor)

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, CursorLocation{Line: 1, Column: 9}, p.buffer.cursor)
	})

	t.Run("Terminate History", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.SetHistory(testHistoryCommands)
		p.keyMapReversed.Insert[Enter] = Terminate
		p.buffer.InsertString("!1")

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, testHistoryCommands[0].Command, p.buffer.String())
		assert.True(t, p.buffer.IsDone())
	})

	t.Run("Terminate Done", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = Terminate
		p.buffer.InsertString("test this thing")

		output := strings.Builder{}
		err := p.handleKeyInsert(termenv.NewOutput(&output), tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, err)
		assert.Equal(t, "test this thing", p.buffer.String())
		assert.True(t, p.buffer.IsDone())
	})

	t.Run("Terminate Insert Newline", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.keyMapReversed.Insert[Enter] = Terminate
		p.buffer.InsertString("test this thing")
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
}
