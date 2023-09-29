package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyntaxHighlighterChroma(t *testing.T) {
	sh, err := SyntaxHighlighterChroma("foo", "bar", "baz")
	assert.Nil(t, sh)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), ErrUnsupportedChromaLanguage)
	assert.Contains(t, err.Error(), "\"foo\"")

	sh, err = SyntaxHighlighterChroma("sql", "bar", "baz")
	assert.NotNil(t, sh)
	assert.Nil(t, err)
}

func TestSyntaxHighlighterSQL(t *testing.T) {
	sh, err := SyntaxHighlighterSQL()
	assert.NotNil(t, sh)
	assert.Nil(t, err)

	assert.Equal(t,
		"\x1b[38;5;81mselect\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;197m*\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;81mfrom\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;231musers\x1b[0m",
		sh(`select * from users`),
	)
}
