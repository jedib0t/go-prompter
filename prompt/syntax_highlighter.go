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
		return nil, fmt.Errorf("%w: %#v (please check Chroma documentation)",
			ErrUnsupportedChromaLanguage, language)
	}
	l = chroma.Coalesce(l)

	f := formatters.Get(formatter)
	if f == nil {
		return nil, fmt.Errorf("%w: %#v (please check Chroma documentation)",
			ErrUnsupportedChromaFormatter, formatter)
	}

	s := styles.Get(style)
	if s == nil {
		return nil, fmt.Errorf("%w: %#v (please check Chroma documentation)",
			ErrUnsupportedChromaStyle, style)
	}

	return func(input string) string {
		if iterator, err := l.Tokenise(nil, input); err == nil {
			rsp := strings.Builder{}
			if err := f.Format(&rsp, s, iterator); err == nil {
				return rsp.String()
			}
		}
		return input
	}, nil
}

// SyntaxHighlighterSQL uses Chroma to return a SQL syntax highlighter.
func SyntaxHighlighterSQL() (SyntaxHighlighter, error) {
	return SyntaxHighlighterChroma("sql", "terminal256", "monokai")
}