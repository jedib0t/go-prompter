package prompt

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func compareLines(t *testing.T, expected, actual []string, msg ...any) {
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
	out := strings.Builder{}

	p := &prompt{}
	err := p.SetKeyMap(KeyMapDefault)
	if err != nil {
		t.Errorf("failed to set up key-map: %v", err)
		t.FailNow()
	}
	p.SetHistoryExecPrefix("!")
	p.SetHistoryListPrefix("!!")
	p.SetInput(os.Stdin)
	p.SetOutput(&out)
	p.SetPrefixer(PrefixText("[" + t.Name() + "] "))
	p.SetRefreshInterval(DefaultRefreshInterval)
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

		p.buffer.InsertString(`select` + ` * from dual`)
		p.updateModel(true)
		expectedLines := []string{
			"[TestPrompt_updateModel/simple_one-liner] select * from dual\x1b[38;5;232;48;5;6m \x1b[0m",
		}
		compareLines(t, expectedLines, p.linesToRender)
	})

	t.Run("simple one-liner with header and footer", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.SetHeader("header")
		p.SetFooter("footer")
		p.updateHeaderAndFooter()

		p.buffer.InsertString(`select` + ` * from dual`)
		p.updateModel(true)
		expectedLines := []string{
			"header",
			"[TestPrompt_updateModel/simple_one-liner_with_header_and_footer] select * from dual\x1b[38;5;232;48;5;6m \x1b[0m",
			"footer",
		}
		compareLines(t, expectedLines, p.linesToRender)
	})

	t.Run("simple one-liner with line-numbers", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.Style().LineNumbers = StyleLineNumbersEnabled
		p.init(ctx)

		p.buffer.InsertString(`select` + ` * from dual`)
		p.updateModel(true)
		expectedLines := []string{
			"[TestPrompt_updateModel/simple_one-liner_with_line-numbers] \x1b[38;5;239;48;5;235m 1 \x1b[0m select * from dual\x1b[38;5;232;48;5;6m \x1b[0m",
		}
		compareLines(t, expectedLines, p.linesToRender)
	})

	t.Run("simple one-liner with line-numbers and short-display-width", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.Style().LineNumbers = StyleLineNumbersEnabled
		p.Style().Dimensions.WidthMin = 95
		p.Style().Dimensions.WidthMax = 95
		p.init(ctx)

		p.buffer.InsertString(`select` + ` * from dual`)
		p.updateModel(false)
		expectedLines := []string{
			"[TestPrompt_updateModel/simple_one-liner_with_line-numbers_and_short-display-width] \x1b[38;5;239;48;5;235m 1 \x1b[0m select ",
			"[TestPrompt_updateModel/simple_one-liner_with_line-numbers_and_short-display-width] \x1b[38;5;239;48;5;235m   \x1b[0m * from ",
			"[TestPrompt_updateModel/simple_one-liner_with_line-numbers_and_short-display-width] \x1b[38;5;239;48;5;235m   \x1b[0m dual",
		}
		compareLines(t, expectedLines, p.linesToRender)
	})

	t.Run("multi-liner with line-numbers and scroll-bar", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.SetAutoCompleter(AutoCompleteSQLKeywords())
		p.SetSyntaxHighlighter(syntaxHighlighter)
		p.Style().LineNumbers = StyleLineNumbersEnabled
		p.Style().LineNumbers.ZeroPrefixed = true
		p.Style().Dimensions.HeightMin = 5
		p.Style().Dimensions.HeightMax = 5
		p.init(ctx)

		testInput := "foo\nbar\nbaz\n"
		p.buffer.InsertString(testInput)
		p.updateModel(true)
		expectedLines := []string{
			"[TestPrompt_updateModel/multi-liner_with_line-numbers_and_scroll-bar] \x1b[38;5;239;48;5;235m 1 \x1b[0m \x1b[38;5;231mfoo\x1b[0m\x1b[38;5;231m",
			"[TestPrompt_updateModel/multi-liner_with_line-numbers_and_scroll-bar] \x1b[38;5;239;48;5;235m 2 \x1b[0m \x1b[0m\x1b[38;5;231mbar\x1b[0m\x1b[38;5;231m",
			"[TestPrompt_updateModel/multi-liner_with_line-numbers_and_scroll-bar] \x1b[38;5;239;48;5;235m 3 \x1b[0m \x1b[0m\x1b[38;5;231mbaz\x1b[0m\x1b[38;5;231m",
			"[TestPrompt_updateModel/multi-liner_with_line-numbers_and_scroll-bar] \x1b[38;5;239;48;5;235m 4 \x1b[0m \x1b[0m\x1b[38;5;232;48;5;6m \x1b[0m",
			"[TestPrompt_updateModel/multi-liner_with_line-numbers_and_scroll-bar] \x1b[38;5;239;48;5;235m   \x1b[0m",
		}
		compareLines(t, expectedLines, p.linesToRender)

		p.buffer.InsertString(testInput)
		p.buffer.InsertString(testInput)
		p.updateModel(true)
		expectedLines = []string{
			"[TestPrompt_updateModel/multi-liner_with_line-numbers_and_scroll-bar] \x1b[38;5;239;48;5;235m 06 \x1b[0m \x1b[0m\x1b[38;5;231mbaz\x1b[0m\x1b[38;5;231m                                         \x1b[38;5;237;48;5;233m░\x1b[0m",
			"[TestPrompt_updateModel/multi-liner_with_line-numbers_and_scroll-bar] \x1b[38;5;239;48;5;235m 07 \x1b[0m \x1b[0m\x1b[38;5;231mfoo\x1b[0m\x1b[38;5;231m                                         \x1b[38;5;237;48;5;233m░\x1b[0m",
			"[TestPrompt_updateModel/multi-liner_with_line-numbers_and_scroll-bar] \x1b[38;5;239;48;5;235m 08 \x1b[0m \x1b[0m\x1b[38;5;231mbar\x1b[0m\x1b[38;5;231m                                         \x1b[38;5;237;48;5;233m░\x1b[0m",
			"[TestPrompt_updateModel/multi-liner_with_line-numbers_and_scroll-bar] \x1b[38;5;239;48;5;235m 09 \x1b[0m \x1b[0m\x1b[38;5;231mbaz\x1b[0m\x1b[38;5;231m                                         \x1b[38;5;237;48;5;233m░\x1b[0m",
			"[TestPrompt_updateModel/multi-liner_with_line-numbers_and_scroll-bar] \x1b[38;5;239;48;5;235m 10 \x1b[0m \x1b[0m\x1b[38;5;232;48;5;6m \x1b[0m                                           \x1b[38;5;237;48;5;233m█\x1b[0m",
		}
		compareLines(t, expectedLines, p.linesToRender)
	})

	t.Run("multi-liner without line-numbers", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.Style().Dimensions.HeightMin = 5
		p.Style().Dimensions.HeightMax = 5
		p.init(ctx)

		testInput := "food\nbard\nbazd\n"
		p.buffer.InsertString(testInput)
		p.updateModel(true)
		expectedLines := []string{
			"[TestPrompt_updateModel/multi-liner_without_line-numbers] food",
			"[TestPrompt_updateModel/multi-liner_without_line-numbers] bard",
			"[TestPrompt_updateModel/multi-liner_without_line-numbers] bazd",
			"[TestPrompt_updateModel/multi-liner_without_line-numbers] \x1b[38;5;232;48;5;6m \x1b[0m",
			"",
		}
		compareLines(t, expectedLines, p.linesToRender)
	})

	t.Run("with auto-complete", func(t *testing.T) {
		p := generateTestPrompt(t, ctx)
		p.SetAutoCompleter(AutoCompleteSQLKeywords())
		p.SetAutoCompleterContextual(AutoCompleteSimple(testSuggestions, true))
		p.Style().LineNumbers = StyleLineNumbersEnabled
		p.Style().LineNumbers.Color = Color{}

		p.buffer.InsertString(`select` + ` * from dual`)
		p.buffer.Insert('\n')
		p.buffer.InsertString(`  where row`)
		p.forceAutoComplete(true)
		p.updateSuggestionsInternal("", "", -1)
		p.updateModel(true)
		expectedLines := []string{
			"[TestPrompt_updateModel/with_auto-complete]  1  select * from dual",
			"[TestPrompt_updateModel/with_auto-complete]  2    where row\x1b[38;5;232;48;5;6m \x1b[0m",
			"[TestPrompt_updateModel/with_auto-complete]            \x1b[38;5;16;48;5;214m row       \x1b[0m\x1b[38;5;16;48;5;208m                \x1b[0m",
			"[TestPrompt_updateModel/with_auto-complete]            \x1b[38;5;16;48;5;45m row_count \x1b[0m\x1b[38;5;0;48;5;39m                \x1b[0m",
			"[TestPrompt_updateModel/with_auto-complete]            \x1b[38;5;16;48;5;45m rownum    \x1b[0m\x1b[38;5;0;48;5;39m Number of Rows \x1b[0m",
			"[TestPrompt_updateModel/with_auto-complete]            \x1b[38;5;16;48;5;45m rows      \x1b[0m\x1b[38;5;0;48;5;39m                \x1b[0m",
		}
		compareLines(t, expectedLines, p.linesToRender)
	})
}
