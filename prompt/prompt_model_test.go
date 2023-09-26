package prompt

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func compareModelLines(t *testing.T, expected, actual []string, msg ...any) {
	assert.Len(t, actual, len(expected))
	assert.Equal(t, expected, actual)

	if fmt.Sprint(expected) != fmt.Sprint(actual) {
		expectedLinesBuilder := strings.Builder{}
		expectedLinesBuilder.WriteString("expectedLines := []string{\n")
		for _, line := range actual {
			expectedLinesBuilder.WriteString(fmt.Sprintf("    %#v,\n", line))
			fmt.Println(line)
		}
		expectedLinesBuilder.WriteString("}\n")
		fmt.Println(expectedLinesBuilder.String())
	} else {
		title := "Actual"
		if len(msg) > 0 {
			title = fmt.Sprint(msg...)
		}
		t.Logf("%s:\n%s", title, strings.Join(actual, "\n"))
	}
}

func generateTestPrompt(t *testing.T, ctx context.Context) *prompt {
	p := &prompt{}
	err := p.SetKeyMap(KeyMapDefault)
	if err != nil {
		t.Errorf("failed to set up key-map: %v", err)
		t.FailNow()
	}
	p.SetHistoryExecPrefix("!")
	p.SetHistoryListPrefix("!!")
	p.SetPrefixer(PrefixText("[" + t.Name() + "] "))
	p.SetRefreshInterval(defaultRefreshInterval)
	p.SetStyle(StyleDefault)
	p.SetTerminationChecker(TerminationCheckerNone())
	p.SetWidthEnforcer(WidthEnforcerDefault)
	p.Style().Cursor.Blink = false
	p.Style().Dimensions.WidthMin = 120
	p.Style().Dimensions.WidthMax = 120
	p.initSync(ctx)
	return p
}

func TestPrompt_updateModel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	syntaxHighlighter, err := SyntaxHighlighterSQL()
	if err != nil {
		t.Errorf("failed to set up syntax-highlighting: %v", err)
		t.FailNow()
	}

	t.Run("simple one-liner", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.SetAutoCompleter(AutoCompleteSQLKeywords())
		p.SetSyntaxHighlighter(syntaxHighlighter)

		p.buffer.InsertString(`select` + ` * from dual`)
		p.updateModel(true)
		expectedLines := []string{
			"[TestPrompt_updateModel/simple_one-liner] \x1b[38;5;81mselect\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;197m*\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;81mfrom\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;231mdual\x1b[0m\x1b[38;5;232;48;5;6m \x1b[0m",
		}
		compareModelLines(t, expectedLines, p.linesToRender)
	})

	t.Run("simple one-liner with line-numbers", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.SetAutoCompleter(AutoCompleteSQLKeywords())
		p.SetSyntaxHighlighter(syntaxHighlighter)
		p.Style().LineNumbers = StyleLineNumbersEnabled
		p.init(ctx)

		p.buffer.InsertString(`select` + ` * from dual`)
		p.updateModel(true)
		expectedLines := []string{
			"[TestPrompt_updateModel/simple_one-liner_with_line-numbers] \x1b[38;5;237;48;5;233m 1 \x1b[0m \x1b[38;5;81mselect\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;197m*\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;81mfrom\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;231mdual\x1b[0m\x1b[38;5;232;48;5;6m \x1b[0m",
		}
		compareModelLines(t, expectedLines, p.linesToRender)
	})

	t.Run("with auto-complete", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.SetAutoCompleter(AutoCompleteSQLKeywords())
		p.SetSyntaxHighlighter(syntaxHighlighter)

		p.buffer.InsertString(`select` + ` * from dual`)
		p.buffer.Insert('\n')
		p.buffer.InsertString(`  where row`)
		p.updateSuggestionsInternal("", "", -1)
		p.updateModel(true)
		expectedLines := []string{
			"[TestPrompt_updateModel/with_auto-complete] \x1b[38;5;81mselect\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;197m*\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;81mfrom\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;231mdual\x1b[0m\x1b[38;5;231m",
			"[TestPrompt_updateModel/with_auto-complete]   \x1b[0m\x1b[38;5;81mwhere\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;81mrow\x1b[0m\x1b[38;5;232;48;5;6m \x1b[0m",
			"[TestPrompt_updateModel/with_auto-complete]        \x1b[38;5;16;48;5;214m row_count \x1b[0m\x1b[38;5;16;48;5;208m                \x1b[0m",
			"[TestPrompt_updateModel/with_auto-complete]        \x1b[38;5;16;48;5;45m rownum    \x1b[0m\x1b[38;5;0;48;5;39m Number of Rows \x1b[0m",
			"[TestPrompt_updateModel/with_auto-complete]        \x1b[38;5;16;48;5;45m rows      \x1b[0m\x1b[38;5;0;48;5;39m                \x1b[0m",
		}
		compareModelLines(t, expectedLines, p.linesToRender)
	})
}
