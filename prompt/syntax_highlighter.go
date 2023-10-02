package prompt

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// SyntaxHighlighter helps do syntax highlighting by using ANSI color codes on
// the prompt text.
type SyntaxHighlighter func(input string) string

// SyntaxHighlighterChroma uses the "github.com/alecthomas/chroma" library to do
// the syntax highlighting. Please refer to the documentation on that project
// for possible values for "language", "formatter" and "style".
func SyntaxHighlighterChroma(language, formatter, style string) (SyntaxHighlighter, error) {
	l := lexers.Get(language)
	if l == nil {
		return nil, fmt.Errorf("%w: %#v (check Chroma lexers documentation)",
			ErrUnsupportedChromaLanguage, language)
	}
	l = chroma.Coalesce(l)
	f := formatters.Get(formatter)
	s := styles.Get(style)

	return func(input string) string {
		output := input
		if iterator, err := l.Tokenise(nil, input); err == nil {
			out := strings.Builder{}
			if err = f.Format(&out, s, iterator); err == nil {
				output = out.String()
			}
		}
		return output
	}, nil
}

// SyntaxHighlighterSQL uses Chroma to return a SQL syntax highlighter.
func SyntaxHighlighterSQL() (SyntaxHighlighter, error) {
	return SyntaxHighlighterChroma("sql", "terminal256", "monokai")
}
