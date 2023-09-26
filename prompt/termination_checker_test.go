package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerminationCheckerNone(t *testing.T) {
	tc := TerminationCheckerNone()

	assert.True(t, tc("foo"))
}

func TestTerminationCheckerSQL(t *testing.T) {
	tc := TerminationCheckerSQL()

	assert.False(t, tc("foo"))
	assert.True(t, tc("/foo"))
	assert.True(t, tc("foo;"))
}
