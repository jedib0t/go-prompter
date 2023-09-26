package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyntaxHighlighterSQL(t *testing.T) {
	sh, err := SyntaxHighlighterSQL()
	assert.NotNil(t, sh)
	assert.Nil(t, err)

	assert.Equal(t,
		"\x1b[38;5;81mselect\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;197m*\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;81mfrom\x1b[0m\x1b[38;5;231m \x1b[0m\x1b[38;5;231musers\x1b[0m",
		sh(`select * from users`),
	)
}
