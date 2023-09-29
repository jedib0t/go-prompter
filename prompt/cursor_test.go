package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCursorLocationString(t *testing.T) {
	cl := CursorLocation{Line: 12, Column: 13}

	assert.Equal(t, "[13, 14]", cl.String())
}
